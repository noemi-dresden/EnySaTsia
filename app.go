package main

import (
	U "EnySaTsia/routes/users"
	R "EnySaTsia/routes/votes"

	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
)

// App -- the application root
func App() *iris.Application {
	app := iris.New()

	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	app.Get("/", func(ctx iris.Context) {
		ctx.JSON("every thing ok")
	})

	s := app.Party("/session", jwtHandler.Serve)
	{
		s.Post("/new", R.NewSession)
		s.Put("/update", R.UpdateSession)
		s.Post("/start", R.StartSession)
		s.Post("/close", R.CloseSession)
		s.Get("/", R.GetAllSession)
		s.Get("/{id: string}", R.GetSession)
	}

	q := app.Party("/question", jwtHandler.Serve)
	{
		q.Post("/new", R.NewQuestion)
		q.Put("/update", R.UpdateQuestion)
		q.Post("/startVote", R.VoteStart)
		q.Post("/closeVote", R.VoteClose)
		q.Post("/vote", R.Vote)
		q.Get("/{question: string}", R.GetQuestion)
		q.Get("/session/{session: string}", R.GetQuestionOfSession)
	}
	u := app.Party("/user")
	{
		u.Post("/register", U.Register)
		u.Post("/login", U.Login)
		u.Post("/changePassword", jwtHandler.Serve, U.ChangePassword)
	}

	return app
}

func main() {
	app := App()
	app.Run(iris.Addr(":8080"))
}
