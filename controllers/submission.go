package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/low4ey/OJ/Golang-backend/database"
	"github.com/low4ey/OJ/Golang-backend/middleware"
	"github.com/low4ey/OJ/Golang-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var SubmissionCollection *mongo.Collection = database.SubmissionData(database.Client, "Submission")

func Submit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Second)
		defer cancel()

		var submission models.Submission
		if err := c.BindJSON(&submission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		testcases, testCaserr := getTestCases("https://serene-fortress-91389-77d1fb95872a.herokuapp.com/api/getTestCase/" + *submission.QuestionId)
		if testCaserr != nil {
			fmt.Println("Error in testcase route")
			c.JSON(http.StatusBadRequest, gin.H{"error": testCaserr.Error()})
			return
		}

		if err := middleware.WriteOutputToFile(testcases); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		submission.SubmitTime = time.Now()

		outcome, status, codeErr := middleware.ExecuteCode(*submission.Code, *submission.Language, testcases)
		if codeErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": codeErr.Error()})
			return
		}
		fmt.Println(outcome, " ", status)
		submission.Status = &status
		submission.LastExecutedIndex = outcome
		submission.Id = primitive.NewObjectID()

		_, err := SubmissionCollection.InsertOne(ctx, submission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert submission"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": submission})
	}
}

func Run() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
		defer cancel()

		var submission models.Submission
		if err := c.BindJSON(&submission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		testCaseURL := "https://serene-fortress-91389-77d1fb95872a.herokuapp.com/api/getTestCase/" + *submission.QuestionId
		sampleTestcase, testCaseErr := getSampleTestCase(testCaseURL)
		if testCaseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": testCaseErr.Error()})
			return
		}

		if err := middleware.WriteOutputToFile(sampleTestcase); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		submission.SubmitTime = time.Now()

		outcome, status, codeErr := middleware.ExecuteCode(*submission.Code, *submission.Language, sampleTestcase)
		if codeErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": codeErr.Error()})
			return
		}

		submission.Status = &status
		submission.LastExecutedIndex = outcome
		submission.Id = primitive.NewObjectID()

		c.JSON(http.StatusCreated, gin.H{"data": submission})
	}
}

func GetAllSub() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := SubmissionCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
		defer cursor.Close(ctx)

		var submissions []models.Submission
		if err := cursor.All(ctx, &submissions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}

		c.JSON(http.StatusOK, submissions)
	}
}

func GetSubByQuestionId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		quesId := c.Param("questionId")
		userId := c.Query("userId")

		// Create the filter
		filter := bson.M{"questionid": quesId}
		if userId != "" {
			filter["userid"] = userId
		}

		// Find options with sort
		findOptions := options.Find().SetSort(bson.D{{Key: "submittime", Value: -1}})

		cursor, err := SubmissionCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
		defer cursor.Close(ctx)

		var submissions []models.Submission
		if err := cursor.All(ctx, &submissions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, submissions)
	}
}

func GetSubByUserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userID := c.Param("userId")

		// Get limit from query string, default to 10 if not provided or invalid
		limitParam := c.Query("limit")
		limit, err := strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			limit = 10 // Default limit if invalid or not provided
		}

		// Define the filter to get submissions for the userID where status is "CORRECT"
		filter := bson.M{
			"userid": userID,
			"status": "CORRECT",
		}

		// Define options to limit and sort the results
		findOptions := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "_id", Value: -1}}) // Sorting by _id descending

		cursor, err := SubmissionCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
		defer cursor.Close(ctx)

		var submissions []models.Submission
		if err := cursor.All(ctx, &submissions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}

		c.JSON(http.StatusOK, submissions)
	}
}
