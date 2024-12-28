package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"testgenerate_backend_user/internal/app"
	"time"
)

type Service interface {
	GetRoles(ctx context.Context) ([]app.Role, error)
	GetUser(ctx context.Context, userName, userRole string) (app.User, error)
	GetUsersRole(ctx context.Context) ([]app.User, error)
	AddUser(ctx context.Context, userAdd app.User) error
	UpdateUser(ctx context.Context, user app.User) error
	DeleteUser(ctx context.Context, userName string) error
}

type userService struct {
	logger *logrus.Logger
}

func NewBasicService(logger *logrus.Logger) Service {
	return userService{
		logger: logger,
	}
}

func NewService(logger *logrus.Logger, requestCount metrics.Counter, requestLatency metrics.Histogram) Service {
	var svc Service
	{
		svc = NewBasicService(logger)
		svc = LoggingMiddleware(logger)(svc)
		svc = InstrumentingMiddleware(requestCount, requestLatency)(svc)
	}
	return svc
}

// ----------------------------------------------------------------------------------------------------------------------
func (u userService) GetRoles(ctx context.Context) ([]app.Role, error) {
	var roles []app.Role
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable",
		app.GetEnv("DB_HOST", "localhost"), app.GetEnvAsInt("DB_PORT", 5432),
		app.GetEnv("DB_USER", "postgres"), app.GetEnv("DB_PASSWORD", "pgpassword"),
		app.GetEnv("DB_NAME", "generate"))
	conn, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		erRet := fmt.Errorf("GetRoles. Unable to connect to database: %v\n", err)
		return roles, erRet
	}
	defer conn.Close(ctx)

	rows, errRows := conn.Query(ctx, `select to_json(t.*)
					from (select id, role_name from user_role) t`)
	if errRows != nil {
		erResp := fmt.Errorf("GetRoles QueryRow: %v\n", errRows)
		return roles, erResp
	}

	for rows.Next() {
		var res string
		errScan := rows.Scan(&res)
		if errScan != nil {
			erRet := fmt.Errorf("GetRoles rows.Scan: %v\n", errScan)
			return roles, erRet
		}
		var result app.Role
		errU := json.Unmarshal([]byte(res), &result)
		if errU != nil {
			erRet := fmt.Errorf("GetRoles json.Unmarshal: %v\n", errU)
			return roles, erRet
		}
		roles = append(roles, result)
	}
	return roles, nil
}
func (u userService) GetUser(ctx context.Context, user, role string) (app.User, error) {
	var userRole app.User
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable",
		app.GetEnv("DB_HOST", "localhost"), app.GetEnvAsInt("DB_PORT", 5432),
		app.GetEnv("DB_USER", "postgres"), app.GetEnv("DB_PASSWORD", "pgpassword"),
		app.GetEnv("DB_NAME", "generate"))
	conn, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		erRet := fmt.Errorf("GetUser. Unable to connect to database: %v\n", err)
		return userRole, erRet
	}
	defer conn.Close(ctx)

	err = conn.QueryRow(context.Background(),
		`select users.user_name, ur.role_name, users.create_time::date 
				from users left join user_role ur on ur.id = users.role 
                where users.user_name = $1`, user).Scan(&userRole.Name, &userRole.Role, userRole.CreateTime)
	if err != nil {
		erRet := fmt.Errorf("GetUser. QueryRow: %v\n", err)
		return userRole, erRet
	}

	return userRole, nil
}
func (u userService) GetUsersRole(ctx context.Context) ([]app.User, error) {
	var users []app.User
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable",
		app.GetEnv("DB_HOST", "localhost"), app.GetEnvAsInt("DB_PORT", 5432),
		app.GetEnv("DB_USER", "postgres"), app.GetEnv("DB_PASSWORD", "pgpassword"),
		app.GetEnv("DB_NAME", "generate"))
	conn, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		erRet := fmt.Errorf("GetUsersRole. Unable to connect to database: %v\n", err)
		return users, erRet
	}
	defer conn.Close(ctx)

	rows, errRows := conn.Query(ctx, `select to_json(t.*)
					from (select users.user_name, ur.role_name,ur.id as role_id, users.create_time::date
							from users left join user_role ur on ur.id = users.role) t`)
	if errRows != nil {
		erResp := fmt.Errorf("GetUsersRole QueryRow: %v\n", errRows)
		return users, erResp
	}

	for rows.Next() {
		var res string
		errScan := rows.Scan(&res)
		if errScan != nil {
			erRet := fmt.Errorf("GetUsersRole rows.Scan: %v\n", errScan)
			return users, erRet
		}
		var result app.User
		errU := json.Unmarshal([]byte(res), &result)
		if errU != nil {
			erRet := fmt.Errorf("GetUsersRole json.Unmarshal: %v\n", errU)
			return users, erRet
		}
		users = append(users, result)
	}
	return users, nil
}
func (u userService) AddUser(ctx context.Context, userAdd app.User) error {
	var errA error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable",
		app.GetEnv("DB_HOST", "localhost"), app.GetEnvAsInt("DB_PORT", 5432),
		app.GetEnv("DB_USER", "postgres"), app.GetEnv("DB_PASSWORD", "pgpassword"),
		app.GetEnv("DB_NAME", "generate"))
	conn, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		erRet := fmt.Errorf("AddUser. Unable to connect to database: %v\n", err)
		return erRet
	}
	defer conn.Close(ctx)

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		errA = fmt.Errorf("AddUser conn.BeginTx %v\n", err)
		return errA
	}
	defer func() {
		if errA != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	//All users add with role == 'user'
	//Next Administrator may change this role
	//SuperAdmins insert trough database
	_, err = tx.Exec(ctx, `insert into users(user_name, role, create_time) values($1, $2)`,
		userAdd.Name, 3, time.Now())
	if err != nil {
		errA = fmt.Errorf("AddUser insert into user: %v\n", err)
		return errA
	}
	return nil
}
func (u userService) UpdateUser(ctx context.Context, user app.User) error {
	var errU error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable",
		app.GetEnv("DB_HOST", "localhost"), app.GetEnvAsInt("DB_PORT", 5432),
		app.GetEnv("DB_USER", "postgres"), app.GetEnv("DB_PASSWORD", "pgpassword"),
		app.GetEnv("DB_NAME", "generate"))
	conn, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		errU = fmt.Errorf("UpdateUser. Unable to connect to database: %v\n", err)
		return errU
	}
	defer conn.Close(ctx)

	_, errU = conn.Exec(ctx, `update users set role = $2, create_time = $3 where user_name = $1`,
		user.Name, user.RoleID, time.Now())
	if errU != nil {
		return fmt.Errorf("UpdateUser conn.Exec: %v\n", errU)
	}
	return nil
}
func (u userService) DeleteUser(ctx context.Context, user string) error {
	var errD error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable",
		app.GetEnv("DB_HOST", "localhost"), app.GetEnvAsInt("DB_PORT", 5432),
		app.GetEnv("DB_USER", "postgres"), app.GetEnv("DB_PASSWORD", "pgpassword"),
		app.GetEnv("DB_NAME", "generate"))
	conn, err := pgx.Connect(ctx, psqlInfo)
	if err != nil {
		errD = fmt.Errorf("DeleteUser. Unable to connect to database: %v\n", err)
		return errD
	}
	defer conn.Close(ctx)

	_, errD = conn.Exec(ctx, `delete from users where user_name = $1`, user)
	if errD != nil {
		return fmt.Errorf("DeleteUser conn.Exec: %v\n", errD)
	}

	return nil
}
