package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bertoxic/graphql2/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *DB) GetJob(id string) *model.JobListing {

	jobCollects := db.client.Database("graphql-jobs").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	defer cancel()
	var jobListing model.JobListing
	err := jobCollects.FindOne(ctx, filter).Decode(&jobListing)
	if err != nil {
		log.Fatal(err)

	}

	return &jobListing
}

func (db *DB) GetJobs() []*model.JobListing {
	jobCollects := db.client.Database("graphql-jobs").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var listofJobs []*model.JobListing
	cursor, err := jobCollects.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(context.TODO(), &listofJobs); err != nil {
		panic(err)
	}
	return listofJobs
}

func (db *DB) CreateJobListing(jobInfo model.CreateJobListingInput) *model.JobListing {

	jobCollects := db.client.Database("graphql-jobs").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	addressx := model.Address{City: *jobInfo.Address.City, Street: *jobInfo.Address.Street, State: *jobInfo.Address.State, Zip: *jobInfo.Address.Zip}
	inserted, err := jobCollects.InsertOne(ctx, bson.M{
		"title":       jobInfo.Title,
		"description": jobInfo.Description,
		"url":         jobInfo.URL,
		"company":     jobInfo.Company,
		"address":     addressx,
	})
	if err != nil {
		log.Fatal(err)
	}
	insertedID := inserted.InsertedID.(primitive.ObjectID).Hex()
	returnjobListing := model.JobListing{ID: insertedID, Title: jobInfo.Title, Company: jobInfo.Company, URL: jobInfo.URL, Description: jobInfo.Description, Address: &addressx}

	return &returnjobListing
}

func (db *DB) UpdateJoblisting(jobID string, jobInfo model.UpdateJobListingInput) *model.JobListing {
	jobCollects := db.client.Database("graphql-jobs").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	updateJobInfo := bson.M{}
    addressx := model.Address{City: *jobInfo.Address.City, Street: *jobInfo.Address.Street, State: *jobInfo.Address.State, Zip: *jobInfo.Address.Zip}
	if jobInfo.Title != nil {
		updateJobInfo["title"] = jobInfo.Title
	}
	if jobInfo.Description != nil {
		updateJobInfo["description"] = jobInfo.Description
	}
	if jobInfo.Company != nil {
		updateJobInfo["company"] = jobInfo.Company
	}
	if jobInfo.URL != nil {
		updateJobInfo["url"] = jobInfo.URL
	}
	if jobInfo.Address != nil {
		updateJobInfo["address"] = addressx
	}
	_id, _ := primitive.ObjectIDFromHex(jobID)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateJobInfo}

	results := jobCollects.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var jobListing model.JobListing
	if err := results.Decode(&jobListing); err != nil {
		log.Fatal(err)
	}
	return &jobListing
}

func (db *DB) DeleteJobListing(jobId string) *model.DeleteJobResponse {
	jobCollects := db.client.Database("graphql-jobs").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	_, err := jobCollects.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	return &model.DeleteJobResponse{DeletedJobID: jobId}
}

func (db *DB) FilterJobs(filter ,field string) ([]*model.JobListing,error){
	jobCollects := db.client.Database("graphql-jobs").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	filterx := bson.M{}
	if field =="_id" {
		convertedID, _ := primitive.ObjectIDFromHex(filter)
		filterx["_id"] = convertedID
	}else{
		filterx = bson.M{fmt.Sprintf("%s",field): filter}
	}
	cux , err := jobCollects.Find(ctx, filterx)
	if err != nil {
		log.Fatal(err)
	}
	var listofJobs []*model.JobListing
	if err = cux.All(context.TODO(),&listofJobs); err !=nil {
		panic(err)
	}
	

	return listofJobs, nil
}
