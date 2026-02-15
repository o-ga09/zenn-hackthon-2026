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
		users.GET("", s.User.List)           // ユーザー一覧取得
		users.GET("/:id", s.User.GetByID)    // IDでユーザー取得
		users.GET("/name", s.User.GetByName) // 名前でユーザー取得
		users.POST("", s.User.Create)        // ユーザー作成
		users.PUT("/:id", s.User.Update)     // ユーザー更新
		users.DELETE("/:id", s.User.Delete)  // ユーザー削除
	}

	// 画像管理API
	images := apiRoot.Group("/media", AuthMiddleware())
	{
		images.GET("", s.Image.List)                          // 画像一覧取得
		images.GET("/:key", s.Image.GetByKey)                 // 画像取得
		images.DELETE("/:key", s.Image.Delete)                // 画像削除
		images.GET("/:id/analytics", s.Image.GetAnalytics)    // 分析結果取得
		images.PUT("/:id/analytics", s.Image.UpdateAnalytics) // 分析結果更新
	}

	// VLog管理API
	vlogs := apiRoot.Group("/vlogs", AuthMiddleware())
	{
		vlogs.GET("", s.VLog.List)                    // VLog一覧取得
		vlogs.GET("/:id", s.VLog.GetByID)             // IDでVLog取得
		vlogs.GET("/:id/stream", s.VLog.StreamStatus) // VLog進捗ストリーミング
		vlogs.DELETE("/:id", s.VLog.Delete)           // VLog削除
	}

	// AIエージェントAPI
	agentGroup := apiRoot.Group("/agent", AuthMiddleware())
	{
		agentGroup.POST("/create-vlog", s.Agent.CreateVLog)                   // VLog作成
		agentGroup.POST("/analyze-media", s.Agent.AnalyzeMedia)               // メディア分析
		agentGroup.GET("/analyze-media/stream", s.Agent.StreamAnalysisStatus) // メディア分析進捗ストリーミング（?ids=id1,id2,id3）
	}

	// 通知API
	notifications := apiRoot.Group("/notifications", AuthMiddleware())
	{
		notifications.GET("", s.Notification.GetNotifications)          // 通知一覧取得
		notifications.PUT("/:id/read", s.Notification.MarkAsRead)       // 通知を既読にする
		notifications.PUT("/read-all", s.Notification.MarkAllAsRead)    // 全通知を既読にする
		notifications.DELETE("/:id", s.Notification.DeleteNotification) // 通知削除
		notifications.DELETE("", s.Notification.DeleteAllNotifications) // 全通知削除
	}

	// 内部タスクAPI（Cloud Tasksからの呼び出し用、自動でIAM認証される）
	internal := s.Engine.Group("/internal")
	{
		internal.POST("/tasks/create-vlog", s.Agent.ProcessVLogTask)
		internal.POST("/tasks/analyze-media", s.Agent.ProcessMediaAnalysisTask)
	}
}
