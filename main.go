package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-gota/gota/dataframe"
	"github.com/minio/minio-go"
)

type RequestJSON struct {
	Input string `json:"minio-input"`
}

const (
	Valid     = 0
	Suggested = 1
	Invalid   = 2
)

func GetEnvVar(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
		return fallback
	}
}

func GetBucketAndPrefix(path string) (string, string) {
	split := strings.Split(path, "/")
	bucket := split[0]
	prefix := strings.Join(split[1:], "/")

	return bucket, prefix
}

func ReadObject(object *minio.Object) (string, error) {
	buffer := new(strings.Builder)
	_, err := io.Copy(buffer, object)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func GetFromS3(client *minio.Client, path string) []int {
	// Get all objects in path
	bucket, prefix := GetBucketAndPrefix(path)
	objects := client.ListObjects(bucket, prefix, true, nil)
	validation := []int{}

	for object := range objects {
		if object.Err != nil {
			println(object.Err)
			continue
		}

		// Get object
		reader, err := client.GetObject(bucket, object.Key, minio.GetObjectOptions{})

		if err != nil {
			println(err)
			continue
		}

		// Read object
		data, err := ReadObject(reader)

		if err != nil {
			println(err)
			continue
		}

		// Build DataFrame and check null values
		df := dataframe.ReadCSV(strings.NewReader(data))

		for _, colName := range df.Names() {
			col := df.Col(colName)
			colType := col.Type()

			// If column has NaN values, it is invalid
			if col.HasNaN() {
				validation = append(validation, Invalid)
			} else {
				if colType == "int" || colType == "float" {
					// If column is not a string, it is valid
					validation = append(validation, Valid)
				} else {
					// If column is a string, it is suggested
					validation = append(validation, Suggested)
				}
			}
		}
	}

	return validation
}

func main() {
	// Connect to MinIO
	host := GetEnvVar("MINIO_HOST", "localhost")
	port := GetEnvVar("MINIO_PORT", "9000")
	endpoint := host + ":" + port
	user := GetEnvVar("MINIO_USER", "minioadmin")
	pass := GetEnvVar("MINIO_USER", "minioadmin")
	client, err := minio.New(endpoint, user, pass, false)

	if err != nil {
		panic(err)
	}

	// Create Gin app and allow CORS
	app := gin.Default()
	app.Use(cors.Default())

	app.POST("", func(c *gin.Context) {
		var json RequestJSON
		err := c.BindJSON(&json)

		if err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		// Get max from validation array
		validation := GetFromS3(client, json.Input)
		max := Valid

		for _, v := range validation {
			if v > max {
				max = v
			}
		}

		if max == Suggested {
			c.JSON(http.StatusOK, gin.H{
				"result": "suggested",
				"reason": "some columns across all files are not number based",
			})
		} else if max == Invalid {
			c.JSON(http.StatusOK, gin.H{
				"result": "invalid",
				"reason": "some columns across all files have NaN values",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": "valid",
				"reason": "all columns across all files are valid",
			})
		}
	})

	app.Run("0.0.0.0:5000")
}
