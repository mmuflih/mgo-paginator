package mgopaginator

import (
	"fmt"
	"math"

	"gopkg.in/mgo.v2"
)

/**
 * Created by Muhammad Muflih Kholidin
 * at 2021-03-06 22:58:55
 * https://github.com/mmuflih
 * muflic.24@gmail.com
 **/

type Paginator struct {
	Query *mgo.Query
	Page  int
	Size  int
	Sort  string
}

func (mgp Paginator) cloneQuery() *mgo.Query {
	nQ := *mgp.Query
	return &nQ
}

func (mgp Paginator) Paginate(items interface{}) *PaginatorResponse {
	c := count(mgp.cloneQuery())
	if mgp.Page == 0 {
		mgp.Page = 1
	}
	if mgp.Size == 0 {
		mgp.Size = 10
	}
	if mgp.Sort != "" {
		mgp.Query = mgp.Query.Sort(mgp.Sort)
	}
	err := mgp.Query.
		Skip((mgp.Page - 1) * mgp.Size).
		Limit(mgp.Size).All(items)
	if err != nil {
		fmt.Println("Error Create Mgo Paginate", err)
	}
	totalPages := int(math.Ceil(float64(c) / float64(mgp.Size)))
	var prevPage, nextPage *int
	var page int = mgp.Page

	if page > 1 {
		np := page - 1
		prevPage = &np
	}
	if page == totalPages {
	} else {
		np := page + 1
		nextPage = &np
	}
	return &PaginatorResponse{
		Data: items,
		Paginate: &PaginatorData{
			Count:      c,
			Page:       page,
			Size:       mgp.Size,
			TotalPages: totalPages,
			NextPage:   nextPage,
			PrevPage:   prevPage,
		},
	}
}

func count(q *mgo.Query) int {
	c, err := q.Count()
	if err != nil {
		return 0
	}
	return c
}

type PaginatorResponse struct {
	Data     interface{}    `json:"data"`
	Paginate *PaginatorData `json:"paginate"`
}

type PaginatorData struct {
	Count      int  `json:"total"`
	Page       int  `json:"page"`
	Size       int  `json:"size"`
	TotalPages int  `json:"total_pages"`
	NextPage   *int `json:"next_page"`
	PrevPage   *int `json:"prev_page"`
}
