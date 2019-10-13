package db

import (
	"fmt"
	"log"
)

// AddGroup creates a group and adds the creating user to it as an admin
func (pg *PGStore) AddGroup(username string, groupname string) error {
	_, err := pg.Query("INSERT INTO groups (group_name) VALUES ($1)",
		groupname)
	if err != nil {
		return fmt.Errorf("AddGroup: %w", err)
	}
	val, err := pg.Query(
		"INSERT INTO member (member, groups, admin) "+
			"select u.id, g.id, true FROM users u INNER JOIN groups g "+
			"ON g.group_name = $1 WHERE username = $2",
		groupname, username)

	log.Println("weird insert returned:", val)
	if err != nil {
		return fmt.Errorf("AddGroup: %w", err)
	}
	return nil
}

// GetGroupsForUser returns all the groups a user can access
func (pg *PGStore) GetGroupsForUser(u User) ([]Group, error) {
	rows, err := pg.Query(
		"SELECT group_name, groups.id from member "+
			"INNER JOIN users ON users.id = member.member "+
			"INNER JOIN groups ON member.groups = groups.id "+
			"where users.id = $1 ORDER BY group_name ASC",
		u.id)
	if err != nil {
		return nil, fmt.Errorf("GetGroupsForUser: %w", err)
	}
	groups := make([]Group, 0)

	for rows.Next() {
		var name, uuid string
		err = rows.Scan(&name, &uuid)
		if err != nil {
			return nil, fmt.Errorf("GetGroupsForUser: Couldn't scan: %w", err)
		}
		groups = append(groups, Group{name, uuid})
	}

	return groups, nil
}

// GetGroupByID returns a group and a list of all its members
func (pg *PGStore) GetGroupByID(id string) (Group, []GroupMember, error) {
	rows, err := pg.Query("SELECT group_name FROM groups WHERE id = $1", id)
	if err != nil {
		return Group{}, nil, err
	}

	if !rows.Next() {
		return Group{}, nil, ErrNotFound
	}

	var name string
	err = rows.Scan(&name)
	if err != nil {
		return Group{}, nil, err
	}

	members := []GroupMember{}
	rows, err = pg.Query("SELECT username,admin FROM member"+
		" LEFT JOIN users ON member = id"+
		" WHERE groups = $1"+
		" ORDER BY username", id)
	if err != nil {
		return Group{}, nil, err
	}

	for rows.Next() {
		member := GroupMember{}
		rows.Scan(&member.Username, &member.Admin)
		members = append(members, member)
	}

	return Group{
		Name: name,
		UUID: id,
	}, members, nil
}
