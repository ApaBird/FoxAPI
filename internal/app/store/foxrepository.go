package store

import (
	"apimod/internal/app/model"
	"database/sql"
	"math/rand"
)

type FoxRepository struct {
	store *Store
}

func (r *FoxRepository) Create(f *model.Fox) (*model.Fox, error) {
	println("Добавлем...")
	if err := r.store.db.QueryRow(
		"INSERT INTO public.foxs (img_url) VALUES ($1) RETURNING id;",
		f.URL,
	).Scan(&f.ID); err != nil {
		return nil, err
	}

	return f, nil
}

func (r *FoxRepository) FindByID(f *model.Fox) (*model.Fox, error) {
	if err := r.store.db.QueryRow(
		"SELECT id, img_url FROM foxs WHERE id = $1",
		f.ID,
	).Scan(&f.ID, &f.URL); err != nil {
		return nil, err
	}

	return f, nil
}

func (r *FoxRepository) Random(f *model.Fox) (*model.Fox, error) {
	var count int
	var rows *sql.Rows

	rows, err := r.store.db.Query("SELECT id FROM public.foxs")
	if err != nil {
		return nil, err
	}

	if err := r.store.db.QueryRow(
		"SELECT count(id) FROM public.foxs", //нужно получить сраз всех id или сразу получить случайного
	).Scan(&count); err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, &StoreError{
			textError: "0 Записей в таблице foxs",
		}
	}

	i := 0
	r_index := rand.Intn(count)

	for rows.Next() {
		rows.Scan(&f.ID)
		i++
		if i == r_index {
			break
		}
	}

	return r.FindByID(f)
}

func (r *FoxRepository) DeleteFoxByID(f *model.Fox) (*model.Fox, error) {
	if err := r.store.db.QueryRow(
		"DELETE FROM public.foxs WHERE id=$1 RETURNING id;",
		f.ID,
	).Scan(&f.ID); err != nil {
		return nil, err
	}

	return f, nil
}

func (r *FoxRepository) UpdateFoxByID(f *model.Fox) (*model.Fox, error) {
	if err := r.store.db.QueryRow(
		"UPDATE public.foxs SET img_url = $2 WHERE id = $1 RETURNING id;",
		f.ID, f.URL,
	).Scan(&f.ID); err != nil {
		return nil, err
	}

	return f, nil
}
