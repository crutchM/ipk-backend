package repository

import (
	"github.com/jmoiron/sqlx"
	"ipk/domain/model"
	"ipk/domain/model/stat"
)

type StatPostgres struct {
	db *sqlx.DB
}

func NewStatPostgres(db *sqlx.DB) *StatPostgres {
	return &StatPostgres{db: db}
}

func (s StatPostgres) getResult(static []stat.Stat) ([]stat.ResponseStat, error) {
	var resp []stat.ResponseStat
	for _, row := range static {
		var item stat.ResponseStat
		var st []stat.TestResult
		err := s.db.Select(&st, "select id, block, question, answer from results where test=$1", row.Id)
		if err != nil {
			return nil, err
		}
		teacher, err := s.getTeacher(row.UserId)
		if err != nil {
			return nil, err
		}
		item.Teacher = teacher
		emp, err := s.getEmployment(row.Employment)
		if err != nil {
			return nil, err
		}
		item.Employment = emp
		var ids []int
		err = s.db.Select(&ids, "select block from results group by block")
		if err != nil {
			return nil, err
		}
		blocks, err := s.getBlocks(ids, row.Id)
		if err != nil {
			return nil, err
		}
		item.Blocks = blocks
		var expert model.Expert
		err = s.db.Get(&expert, "select * from expert where id=$1", row.Expert)
		if err != nil {
			return nil, err
		}
		item.Expert = expert
		item.LessonDate = row.LessonDate
		item.AnketDate = row.AnketDate
		resp = append(resp, item)
	}

	return resp, nil
}
func (s StatPostgres) getBlocks(ids []int, statId int) ([]model.Block, error) {
	var result []model.Block
	var st []stat.TestResult
	if err := s.db.Select(&st, "select id, block, question, answer from results where test=$1", statId); err != nil {
		return nil, err
	}
	for _, id := range ids {
		var block model.Block
		if err := s.db.Get(&block, "select * from block where id=$1", id); err != nil {
			return nil, err
		}
		var tmp []int
		if err := s.db.Select(&tmp, "select question_id from blockQuestions where block_id=$1", block.Id); err != nil {
			return nil, err
		}
		for _, ques := range tmp {
			var question model.Question
			if err := s.db.Get(&question, "select * from question where id=$1", ques); err != nil {
				return nil, err
			}
			var answer int
			for _, v := range st {
				if v.Question == question.Id {
					answer |= v.Answer
				}
			}
			question.Answer = answer
			block.Questions = append(block.Questions, question)
		}
		result = append(result, block)
	}
	return result, nil
}

func (s StatPostgres) getEmployment(id int) (string, error) {
	var emp string
	err := s.db.Get(&emp, "select name from lessontype where id=$1", id)
	if err != nil {
		return "", err
	}
	return emp, nil
}

func (s StatPostgres) getTeacher(id string) (model.User, error) {
	var teacher model.User
	err := s.db.Get(&teacher, "select id, fullname from users where id=$1", id)
	if err != nil {
		return model.User{}, err
	}
	return teacher, nil
}

func (s StatPostgres) GetStat(chair int) ([]stat.ResponseStat, error) {
	var static []stat.Stat
	err := s.db.Select(&static, "select * from stat where chair=$1", chair)
	if err != nil {
		return nil, err
	}
	return s.getResult(static)
}

func (s StatPostgres) GetStatByTeacher(id string) ([]stat.ResponseStat, error) {
	var static []stat.Stat
	err := s.db.Select(&static, "select * from stat where useri=$1", id)
	if err != nil {
		return nil, err
	}
	return s.getResult(static)
}

func (s StatPostgres) AddRow(stat stat.Stat) (int, error) {
	var id int
	row := s.db.QueryRow("insert into stat(useri, post, chair, employment, expert, lessondate, anketdate) values ($1, $2, $3, $4, $5, $6, $7) returning id",
		stat.UserId, stat.PostId, stat.ChairId, stat.Employment, stat.Expert, stat.LessonDate, stat.AnketDate)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (s StatPostgres) AddResult(result []model.Block, rowId int) error {
	for _, v := range result {
		for _, val := range v.Questions {
			var id int
			row := s.db.QueryRow("insert into results(test, block, question, answer) values ($1,$2,$3,$4) returning id", rowId, v.Id, val.Id, val.Answer)
			if err := row.Scan(&id); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s StatPostgres) RemoveUser(id string) error {
	var stats []int
	err := s.db.Select(&stats, "select id from stat where useri=$1", id)
	if err != nil {
		return err
	}
	if stats == nil {
		s.Remove(id)
		return err
	} else {
		for _, v := range stats {
			s.db.QueryRow("DELETE from results where test=$1", v)
			s.db.QueryRow("DELETE from stat where id=$1", v)
		}
		s.Remove(id)
	}
	return nil
}

func (s StatPostgres) Remove(id string) {
	s.db.QueryRow("delete from users where id=$1", id)
}
