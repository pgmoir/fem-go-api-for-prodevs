package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v4"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	// run migrations
	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migratring test db failed: %v", err)
	}

	_, err = db.Exec(`TRUNCATE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("trauncate of db failed: %v", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		watErr  bool
	}{{
		name: "valid workout",
		workout: &Workout{
			Title:           "push day",
			Description:     "upper body day",
			DurationMinutes: 60,
			CaloriesBurned:  200,
			Entries: []WorkoutEntry{
				{
					ExerciseName: "Bench press",
					Sets:         3,
					Reps:         IntPtr(10),
					Weight:       FloatPtr(135.5),
					Notes:        "Warm up properly",
					OrderIndex:   1,
				},
			},
		},
		watErr: false,
	}, {
		name: "valid with invalid entries",
		workout: &Workout{
			Title:           "full bopdy day",
			Description:     "complete body",
			DurationMinutes: 90,
			CaloriesBurned:  500,
			Entries: []WorkoutEntry{
				{
					ExerciseName: "Plank",
					Sets:         3,
					Reps:         IntPtr(60),
					Notes:        "Keep focus",
					OrderIndex:   1,
				},
				{
					ExerciseName:    "Squats",
					Sets:            4,
					Reps:            IntPtr(12),
					DurationSeconds: IntPtr(60),
					Weight:          FloatPtr(185.0),
					Notes:           "Full depth",
					OrderIndex:      2,
				},
			},
		},
		watErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.watErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)

			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)

			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(createdWorkout.Entries), len(retrieved.Entries))

			for i := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(i float64) *float64 {
	return &i
}
