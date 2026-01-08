package server

func (s *Server) SetupApplicationRoute() {
	apiRoot := s.Engine.Group("/api")

	users := apiRoot.Group("/users")
	{
		users.GET("", s.User.List)          // ユーザー一覧取得
		users.GET("/:id", s.User.GetByID)   // IDでユーザー取得
		users.POST("", s.User.Create)       // ユーザー作成
		users.PUT("/:id", s.User.Update)    // ユーザー更新
		users.DELETE("/:id", s.User.Delete) // ユーザー削除
	}
}
