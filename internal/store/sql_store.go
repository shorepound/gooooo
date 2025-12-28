package store

import (
    "database/sql"
    "errors"
)

type PostgresStore struct {
    db *sql.DB
}

func NewPostgres(db *sql.DB) Backend {
    return &PostgresStore{db: db}
}

func (p *PostgresStore) Create(it Item) Item {
    var id int64
    err := p.db.QueryRow(`INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id`, it.Name, it.Description).Scan(&id)
    if err != nil {
        return Item{}
    }
    it.ID = id
    return it
}

func (p *PostgresStore) List() []Item {
    rows, err := p.db.Query(`SELECT id, name, description FROM items`)
    if err != nil {
        return nil
    }
    defer rows.Close()
    var res []Item
    for rows.Next() {
        var it Item
        _ = rows.Scan(&it.ID, &it.Name, &it.Description)
        res = append(res, it)
    }
    return res
}

func (p *PostgresStore) Get(id int64) (Item, bool) {
    var it Item
    err := p.db.QueryRow(`SELECT id, name, description FROM items WHERE id=$1`, id).Scan(&it.ID, &it.Name, &it.Description)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return Item{}, false
        }
        return Item{}, false
    }
    return it, true
}

func (p *PostgresStore) Update(id int64, in Item) (Item, error) {
    res, err := p.db.Exec(`UPDATE items SET name=$1, description=$2 WHERE id=$3`, in.Name, in.Description, id)
    if err != nil {
        return Item{}, err
    }
    n, err := res.RowsAffected()
    if err != nil {
        return Item{}, err
    }
    if n == 0 {
        return Item{}, errors.New("not found")
    }
    in.ID = id
    return in, nil
}

func (p *PostgresStore) Delete(id int64) bool {
    res, err := p.db.Exec(`DELETE FROM items WHERE id=$1`, id)
    if err != nil {
        return false
    }
    n, err := res.RowsAffected()
    if err != nil {
        return false
    }
    return n > 0
}
