package cmd

import (
	"encoding/binary"
	bolt "go.etcd.io/bbolt"
)

var db *bolt.DB
var taskBucket = []byte("tasks")

type Task struct {
	Key   int
	Value string
}

func DbInit() error {
	var err error
	db, err = bolt.Open("tasks.db", 0666, nil)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		if err != nil {
			return err
		}
		return nil
	})	
}

func CreateTask(title string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id, _ := b.NextSequence()
		return b.Put(itob(int(id)), []byte(title))
	})
}

func TaskList() ([]Task, error) {
	var tasks []Task

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{btoi(k), string(v)})
		}
		return nil
	})

	if err!=nil {
		return nil, err
	}

	return tasks, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

