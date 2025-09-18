package db

import (
	"context"
	"errors"
	"time"

	env "GOLANG_SERVER/components/env"
	"GOLANG_SERVER/components/schema"

	"go.mongodb.org/mongo-driver/bson"
)

func GetDeviceAddress(userID string) ([]schema.GetDevice, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	// เชื่อมต่อกับ MongoDB
	collection := client.Database(env.GetEnv("MONGO_DB")).Collection(env.GetEnv("MONGO_DEVICECOLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ค้นหาเอกสารทั้งหมดที่ตรงกับ userID
	cursor, err := collection.Find(ctx, bson.M{"userID": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// สร้าง slice สำหรับเก็บผลลัพธ์
	var devices []schema.GetDevice
	for cursor.Next(ctx) {
		var device schema.GetDevice
		if err := cursor.Decode(&device); err != nil {
			return nil, err
		}

		// เพิ่มข้อมูลลงใน slice
		devices = append(devices, device)
	}

	// ตรวจสอบข้อผิดพลาดจาก cursor
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}
