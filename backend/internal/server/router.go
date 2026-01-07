package server

func (s *Server) SetupApplicationRoute() {
	users := s.Engine.Group("/users")
	{
		users.GET("/", s.User.GET)
		users.POST("/", s.User.Create)
		users.PUT("/:id", s.User.Update)
		users.DELETE("/:id", s.User.Delete)
	}
}
