package server

func (s *Server) SetupApplicationRoute() {
	apiRoot := s.Engine.Group("/api")
	// 認証API
	auth := apiRoot.Group("/auth")
	{
		auth.POST("/session", s.Auth.SignUp)                // セッション作成
		auth.DELETE("/session", s.Auth.SignOut)             // セッション削除
		auth.GET("/user", s.Auth.GetUser, AuthMiddleware()) // ログインユーザー情報取得（認証必須）
	}

	// ユーザー管理API
	users := apiRoot.Group("/users", AuthMiddleware())
	{
		users.GET("", s.User.List)          // ユーザー一覧取得
		users.GET("/:id", s.User.GetByID)   // IDでユーザー取得
		users.POST("", s.User.Create)       // ユーザー作成
		users.PUT("/:id", s.User.Update)    // ユーザー更新
		users.DELETE("/:id", s.User.Delete) // ユーザー削除
	}
}
