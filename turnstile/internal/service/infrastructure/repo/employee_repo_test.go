package repo

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"regexp"
	"testing"
	"turnstile/internal/models"
	"turnstile/pkg/logging"
	"turnstile/pkg/sqlite"
)

//go test -bench=BenchmarkEmployeeRepo_Save5000 -benchmem -benchtime=10x
//go test -bench=BenchmarkEmployeeRepo_SaveSlice5000 -benchmem -benchtime=10x

const (
	debugLevel = "debug"
)

func BenchmarkEmployeeRepo_Save5000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		saveEmployee(5000, b)
	}
}

func BenchmarkEmployeeRepo_SaveSlice5000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		saveSlice(5000, b)
	}
}

func TestEmployeeRepo_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := logging.GetLogger(debugLevel)

	r := NewEmployeeRepo(logger, dbx, "employee")
	type args struct {
		data models.Employee
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

				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO employee")).
					WithArgs(args.data.CardNumber, args.data.EmployeeID, args.data.Rv, args.data.IsDeleted).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			input: args{
				data: models.Employee{
					CardNumber: 414141,
					EmployeeID: 2222,
					Rv:         27493955359,
					IsDeleted:  true,
				},
			},
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			err = r.Save(tt.input.data)
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

func TestEmployeeRepo_GetEmployeeByCard(t *testing.T) {
	db, mock, err := sqlmock.New()
	dbx := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := logging.GetLogger(debugLevel)
	r := NewEmployeeRepo(logger, dbx, "employee")
	type args struct {
		card uint64
	}

	type mockBehavior func(args args)
	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    models.Employee
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"card_number", "employee_id", "rv", "isdeleted"}).
					AddRow(414141, 2222, 27493955359, true)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM employee")).
					WithArgs(args.card).WillReturnRows(rows)
			},
			input: args{
				card: 414141,
			},
			want: models.Employee{
				CardNumber: 414141, EmployeeID: 2222, Rv: 27493955359, IsDeleted: true,
			},
			wantErr: false,
		},
		{
			name: "No Records",
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"card_number", "employee_id", "rv", "isdeleted"})
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM employee")).
					WithArgs(args.card).WillReturnRows(rows)
			},
			input: args{
				card: 414141,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.GetEmployeeByCard(tt.input.card)
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

func saveSlice(x int, b *testing.B) {
	r := newEmployeeRepo()

	employees := make([]models.Employee, 0)

	for i := 0; i < x; i++ {
		employees = append(employees, generateNewEmployee())
	}

	b.ResetTimer()
	for j := 0; j <= x; j++ {
		err := r.SaveSlice(employees)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func saveEmployee(x int, b *testing.B) {
	r := newEmployeeRepo()

	employee := generateNewEmployee()

	b.ResetTimer()
	for j := 0; j <= x; j++ {
		err := r.Save(employee)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateNewEmployee() models.Employee {
	var maxValue int64 = 70000

	return models.Employee{
		CardNumber: uint64(rand.Int63n(maxValue)),
		EmployeeID: uint64(rand.Int63n(maxValue)),
		Rv:         uint64(rand.Int63n(100 * maxValue)),
		IsDeleted:  false,
	}
}

func newEmployeeRepo() *EmployeeRepo {
	logger := logging.GetLogger(debugLevel)

	db, err := sqlite.New(sqlite.Config{
		FileName: "test.db",
	})
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	return NewEmployeeRepo(logger, db, "employee")
}
