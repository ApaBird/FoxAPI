package apiserver

import (
	"apimod/internal/app/model"
	"apimod/internal/app/store"
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type APIserver struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *APIserver {
	return &APIserver{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIserver) Strat() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	if err := s.configreStore(); err != nil {
		return err
	}

	s.configureRoute()

	s.logger.Info("API сервер запущен")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIserver) configreStore() error {
	s.logger.Info("Подключение к БД...")
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		s.logger.Error("Упс... не удалось подключиться к БД")
		return err
	}

	s.store = st

	s.logger.Info("БД подключена")
	return nil
}

func (s *APIserver) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *APIserver) configureRoute() {
	s.router.HandleFunc("/", s.randomFox())
	s.router.HandleFunc("/getfox", s.changeFox(s.store.Fox().FindByID))
	s.router.HandleFunc("/addfox", s.changeFox(s.store.Fox().Create))
	s.router.HandleFunc("/deletefox", s.changeFox(s.store.Fox().DeleteFoxByID))
	s.router.HandleFunc("/updatefox", s.changeFox(s.store.Fox().UpdateFoxByID))
}

func (s *APIserver) randomFox() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fox *model.Fox = &model.Fox{}
		s.logger.Info("Получаем случайную Лису...")

		if _, err := s.store.Fox().Random(fox); err != nil {
			io.WriteString(w, err.Error())
			return
		}

		s.respondHTML(w, r, http.StatusAccepted, fox, "index.html")
	}
}

func (s *APIserver) changeFox(action func(*model.Fox) (*model.Fox, error)) http.HandlerFunc {
	type request struct {
		Url string `json:"url"`
		ID  int    `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Error("Ошибка данных")
			s.logger.Error(err)
			s.logger.Info(r.Header)
			s.logger.Info(r.Host)
			s.logger.Info(r.Method)

			s.logger.Info(r.FormValue("url"))
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		f := &model.Fox{
			URL: req.Url,
			ID:  req.ID,
		}

		if _, err := action(f); err != nil {
			s.logger.Error("Ошибка с БД")
			s.logger.Error(err)
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusAccepted, f)

	}

}

func (s *APIserver) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *APIserver) errorHTML(w http.ResponseWriter, r *http.Request, code int, err error) {
	t, e := template.ParseFiles("template/Error.html")
	if e != nil {
		s.logger.Error(e)
		return
	}

	s.logger.Error(err)
	t.Execute(w, err)
}

func (s *APIserver) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *APIserver) respondHTML(w http.ResponseWriter, r *http.Request, code int, data interface{}, html string) {
	w.WriteHeader(code)
	if data != nil {
		temp, err := template.ParseFiles("template/" + html)
		if err != nil {
			t, _ := template.ParseFiles("template/Error.html")

			s.logger.Error(err)
			t.Execute(w, err)
		}

		temp.Execute(w, data)
	}
}
