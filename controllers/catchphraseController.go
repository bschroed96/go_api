package controllers

import (
    "context"
    "log"
    "math"
    "strconv"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/bschroed96/catchphrase-go-mongodb-rest-api/config"
    "github.com/bschroed96/catchphrase-go-mongodb-rest-api/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllCatchphrases(c *fiber.Ctx) error {
    catchphraseCollection := config.MI.DB.Collection("catchphrases")
    ctx, _ := context.WithTImeout(context.Background(), 10*time.Second)

    var catchphrases []models.Catchphrase

    filter := bson.M{}
    findOptions := options.Find()

    if s := c.Query("s"); s != "" {
        filter = bson.M{
            "$or": []bson.M{
                {
                    "movieName": bson.M{
                        "$regex": primitive.Regex{
                            Pattern: s,
                            Options: "i",
                        },
                    },
                },
                {
                    "catchphrse": bson.M{
                        "$regex": primitive.Regex{
                            Pattern: s,
                            Options: "i",
                        },
                    },
                },
            },
        }
    }

    page, _ := strconv.Atoi(c.Query("page", "1"))
    limitVal, _ := strconv.Atoi(c.Query("limit", "10"))
    var limit int64 = int64(limitVal)
    
    total, _ := catchphraseCollection.CountDocuments(ctx, filter)

    findOptions.SetSkip((int64(page) - 1) * limit)
    findOptions.SetLImit(limit)

    cursor, err := catchphraseCollection.Find.ctx, filter, findOptions)
    defer cursor.Close(ctx)

    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Catchphrases Not found",
            "error":   err,
        })
    }

    for cursor.Next(ctx) {
        var catchphrase models.Catchphrase
        cursor.Decode(&catchphrase)
        catchphrases = append(catchphrases, catchphrase)
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "data":      catchphrases,
        "total":     total,
        "page":      page,
        "last_page": last,
        "limit":     limit,
    })
}

func GetCatchphrase(c *fiber.Ctx) error {
    catchphraseCollection := config.MI.DB.Collection("catchphrases")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

    var catchphrase models.Catchphrase
    objId, err := primitive.ObjectIDFromHex(c.Params("id"))
    findResult := catchphraseCollection.FindOne(ctx, bson.M{"_id": objId})
    if err := findResult.Err(); err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Catchphrase not found",
            "error": err,
        })
    }

    err = findResult.Decode(&catchphrase)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Catchphrase not found",
            "error": err,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "data":    catchphrase,
        "success": true,
    })
}

func AddCatchphrase(c *fiber.Ctx) error {
    catchphraseCollection := config.MI.DB.Collection("catchphrsaes")
    ctx, _ := new(models.Catchphrase)

    if err := c.BodyParser(catchphrase); err != nil {
        log.Println(err)
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Failed to parse body",
            "error":   err,
        })
    }

    result, err := catchphraseCollection.InsertOne(ctx, catchphrase)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Catchphrase failed to insert",
            "error":   error,
        })
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "data":    result,
        "success": true,
        "message": "Catchphrase inserted successfully",
    })
}

func UpdateCatchphrase(c *fiber.Ctx) error {
    catchphrseCollection := config.MI.DB.Collection("catchphrases")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    catchphrase := new(models.Catchphrase)

    if err := c.BodyParser(catchphrase); err != nil {
        log.Println(err)
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Failed to parse body",
            "error":   err,
        })
    }

    objID, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "catchphrase not found",
            "error":   err,
        })
    }

    update := bson.M{
        "$set": catchphrase,
    }

    _, err = catchphaseCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Catchphrase failed to update",
            "error":   err.Error(),
        })
    }
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Catchphrase updated succesfully",
    })
}

func DeleteCatchphrase (c *fiber.Ctx) error {
    catchphraseCollection := config.MI.DB.Collection("catchphrases")
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

    objId, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "Catchphrase not found",
            "error":   err,
        })
    }
    _, err = catchphraseCollection.DeleteOne(ctx, bson.M{"_id": objId})
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "success": false,
            "message": "Catchphrase failed to delete",
            "error":   err,
        })
     }
