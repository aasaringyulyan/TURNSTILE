package repo

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log"
	"regexp"
	"testing"
	"turnstile/internal/models"
	"turnstile/pkg/logging"
)

func TestLogRepo_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := logging.GetLogger(debugLevel)
	r := NewLogRepo(logger, dbx, "log")

	type args struct {
		log models.PassageLogForApi
	}

	type mockBehavior func(args args)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(args args) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO log")).
					WithArgs(args.log.TurnstileID, args.log.EmployeeID, args.log.CardID, args.log.Direction, args.log.DateTime).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			input: args{
				log: models.PassageLogForApi{
					TurnstileID: 123,
					EmployeeID:  123,
					CardID:      123,
					Direction:   1,
					DateTime:    "123",
				},
			},
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			err = r.Save(tt.input.log)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, 1)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLogRepo_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := logging.GetLogger(debugLevel)
	r := NewLogRepo(logger, dbx, "log")

	type mockBehavior func()

	tests := []struct {
		name    string
		mock    mockBehavior
		want    models.PassageLogsForApi
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"turnstile_id", "employee_id", "card", "direction", "dt"}).
					AddRow(1, 2, 3, 4, "1234")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM log")).WillReturnRows(rows)
			},
			want: models.PassageLogsForApi{
				Logs: []models.PassageLogForApi{
					{1, 2, 3, 4, "1234"},
				},
			},
			wantErr: false,
		},
		{
			name: "No Records",
			mock: func() {
				rows := sqlmock.NewRows([]string{"turnstile_id", "employee_id", "card", "direction", "dt"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM log")).WillReturnRows(rows)
			},
			want: models.PassageLogsForApi{
				Logs: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAll()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLogRepo_DeleteAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := logging.GetLogger(debugLevel)
	r := NewLogRepo(logger, dbx, "log")

	type mockBehavior func()

	tests := []struct {
		name    string
		mock    mockBehavior
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec(regexp.QuoteMeta("DELETE FROM log")).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err = r.DeleteAll()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, 1)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
