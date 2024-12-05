/*
 * CR3 Client Dashboard
 * Package: gitlab.hexacode.org/cr3/client-dashboard
 * Maintainer: Azzis Arswendo <azzis@hexacode.org>
 * 						 Dean Ramhan <dean@hexacode.org>
 *
 * Copyright (C) 2024 Hexacode Teknologi Indonesia
 * All Rights Reserved
 */

package main

import (
	"fmt"
	"strconv"
	"strings"

	U "gitlab.hexacode.org/go-libs/hctypes"
	"gitlab.hexacode.org/go-libs/microservice/config"
	"gitlab.hexacode.org/go-libs/strapi"
)

func CreateEmployee(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	employee_name, ok := ctx.Params.GetString("employee_name")

	if !ok {
		return ctx.Error(400, "invalid parameters employee_name")
	}

	email, ok := ctx.Params.GetString("email")

	if !ok {
		return ctx.Error(400, "invalid parameters email")
	}

	phone, ok := ctx.Params.GetString("phone")

	if !ok {
		return ctx.Error(400, "invalid parameters phone")
	}

	address, ok := ctx.Params.GetString("address")

	if !ok {
		return ctx.Error(400, "invalid parameters address")
	}

	city, ok := ctx.Params.GetString("city")

	if !ok {
		return ctx.Error(400, "invalid parameters city")
	}

	state, ok := ctx.Params.GetString("state")

	if !ok {
		return ctx.Error(400, "invalid parameters state")
	}

	zip_code, ok := ctx.Params.GetString("zip_code")

	if !ok {
		return ctx.Error(400, "invalid parameters zip_code")
	}

	country, ok := ctx.Params.GetString("country")

	if !ok {
		return ctx.Error(400, "invalid parameters	country")
	}

	date_of_birth, ok := ctx.Params.GetString("date_of_birth")

	if !ok {
		return ctx.Error(400, "invalid parameters date_of_birth")
	}

	job_level, ok := ctx.Params.GetString("job_level")

	if !ok {
		return ctx.Error(400, "invalid parameters job_level")
	}

	job_position, ok := ctx.Params.GetString("job_position")

	if !ok {
		return ctx.Error(400, "invalid parameters job_position")
	}

	branch, ok := ctx.Params.GetString("branch")

	if !ok {
		return ctx.Error(400, "invalid parameters branch")
	}

	employment_status, ok := ctx.Params.GetString("employment_status")

	if !ok {
		return ctx.Error(400, "invalid parameters employment_status")
	}

	start_date, ok := ctx.Params.GetString("start_date")

	if !ok {
		return ctx.Error(400, "invalid parameters start_date")
	}

	end_date, ok := ctx.Params.GetString("end_date")

	if !ok {
		end_date = ""
	}

	dataSend := U.Dict{
		"employee_name":     employee_name,
		"email":             email,
		"phone":             phone,
		"address":           address,
		"city":              city,
		"state":             state,
		"zip_code":          zip_code,
		"country":           country,
		"date_of_birth":     date_of_birth,
		"job_level":         job_level,
		"job_position":      job_position,
		"branch":            branch,
		"employment_status": employment_status,
		"start_date":        start_date,
	}

	if end_date != "" {
		dataSend["end_date"] = end_date
	}

	status, employee, err := Strapi.Create("employee-table")

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(employee)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	status, checkPerson, err := Strapi.Find("person-table", U.Dict{
		"filters": U.Dict{
			"email": U.Dict{
				"$eq": email,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		fmt.Println(checkPerson)
		fmt.Println("ayam")
		return ctx.Error(500, "Internal Server Error")
	}

	data, ok := checkPerson.GetList("data")

	if !ok {
		fmt.Println("ayam 1")
		return ctx.Error(500, "Internal Server Error")
	}

	if len(data) > 0 {
		person, ok := data.GetDict(0)

		if !ok {
			fmt.Println("ayam 2")
			return ctx.Error(500, "Internal Server Error")
		}

		idPerson, ok := person.GetFloat64("id")

		if !ok {
			fmt.Println("ayam 3")
			return ctx.Error(500, "Internal Server Error")
		}

		id, ok := employee["data"].(U.Dict).GetFloat64("id")

		if !ok {
			fmt.Println("ayam 4")
			return ctx.Error(500, "Internal Server Error")
		}

		status, updateEmployee, err := Strapi.Update("employee-table", int64(id), U.Dict{
			"person": idPerson,
		})

		if err != nil {
			fmt.Println("ayam 5")
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println("ayam 6")
			fmt.Println(updateEmployee)
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, employee, nil
	} else {
		fmt.Print("employee")
		fmt.Print(employee)
		id, ok := employee["data"].(U.Dict).GetFloat64("id")

		if !ok {
			fmt.Println("ayam 7")
			return ctx.Error(500, "Internal Server Error")
		}

		status, createPerson, err := Strapi.Create("person-table", U.Dict{
			"fullname":      employee_name,
			"email":         email,
			"phone":         phone,
			"address":       address,
			"city":          city,
			"state":         state,
			"zip_code":      zip_code,
			"country":       country,
			"date_of_birth": date_of_birth,
			"is_active":     "Yes",
			"employee":      id,
		})

		if err != nil {
			fmt.Println(err)
			return ctx.Error(500, "Internal Server Error")
		}

		if status != 200 {
			fmt.Println(createPerson)
			fmt.Println(status)
			return ctx.Error(500, "Internal Server Error")
		}

		return 200, employee, nil
	}
}

func GetEmployees(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	page, ok := ctx.Params.GetFloat64("page")
	if !ok {
		fmt.Println("invalid parameters")
		page = 1
	}

	page_size, ok := ctx.Params.GetFloat64("page_size")
	if !ok {
		fmt.Println("invalid parameters")
		page_size = 10
	}

	status, employees, err := Strapi.Find("employee-table", U.Dict{
		"pagination": U.Dict{
			"page":     page,
			"pageSize": page_size,
		},
		"sort": "id:desc",
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(employees)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, employees, nil
}

func GetEmployee(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	employeeId, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(employeeId))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, employee, err := Strapi.Find("employee-table", U.Dict{
		"id": id,
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(employee)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, employee, nil
}

func UpdateEmployee(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	employeeId, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(employeeId))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	employee_name, ok := ctx.Params.GetString("employee_name")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	email, ok := ctx.Params.GetString("email")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	phone, ok := ctx.Params.GetString("phone")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	address, ok := ctx.Params.GetString("address")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	city, ok := ctx.Params.GetString("city")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	state, ok := ctx.Params.GetString("state")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	country, ok := ctx.Params.GetString("country")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	date_of_birth, ok := ctx.Params.GetString("date_of_birth")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	zip_code, ok := ctx.Params.GetString("zip_code")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	job_level, ok := ctx.Params.GetString("job_level")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	job_position, ok := ctx.Params.GetString("job_position")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	branch, ok := ctx.Params.GetString("branch")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	employment_status, ok := ctx.Params.GetString("employment_status")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	start_date, ok := ctx.Params.GetString("start_date")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	end_date, ok := ctx.Params.GetString("end_date")

	if !ok {
		end_date = ""
	}

	is_active, ok := ctx.Params.GetString("is_active")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	dataSend := U.Dict{
		"employee_name":     employee_name,
		"email":             email,
		"phone":             phone,
		"address":           address,
		"city":              city,
		"state":             state,
		"country":           country,
		"date_of_birth":     date_of_birth,
		"zip_code":          zip_code,
		"job_level":         job_level,
		"job_position":      job_position,
		"branch":            branch,
		"employment_status": employment_status,
		"start_date":        start_date,
		"is_active":         is_active,
	}

	if end_date != "" {
		dataSend["end_date"] = end_date
	}

	status, employee, err := Strapi.Update("employee-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(employee)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, employee, nil
}

func DeleteEmployee(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {

	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	employeeId, ok := ctx.URLParams.GetString("id")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	id, err := strconv.Atoi(string(employeeId))

	if err != nil {
		return ctx.Error(400, "invalid parameters")
	}

	status, employee, err := Strapi.Delete("employee-table", int64(id))

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(employee)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return status, employee, nil
}

func GetCountEmployeePerMonth(ctx *config.ConfigContext) (int, U.Dict, *config.ReturnHTTPError) {
	bearer, bearer_ok := ctx.Headers.GetString("Authorization")

	if !bearer_ok {
		return ctx.Error(401, "Unauthorized")
	}

	Strapi := strapi.NewStrapiWithToken(STRAPI_URL, strings.Replace(string(bearer), "Bearer ", "", 1))

	status, me, err := Strapi.UsersMe()
	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	fmt.Println(me)

	earlyDate, ok := ctx.Params.GetString("earlyDate")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	lastDate, ok := ctx.Params.GetString("lastDate")

	if !ok {
		return ctx.Error(400, "invalid parameters")
	}

	status, employee, err := Strapi.Find("employee-table", U.Dict{
		"filters": U.Dict{
			"createdAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
		},
	})

	fmt.Println(employee)

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	status, checkinEmployee, err := Strapi.Find("attendance-table", U.Dict{
		"filters": U.Dict{
			"createdAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
			"payload": U.Dict{
				"$containsi": "Checkin",
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("checkinEmployee :", checkinEmployee)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	status, absentEmployee, err := Strapi.Find("attendance-table", U.Dict{
		"filters": U.Dict{
			"createdAt": U.Dict{
				"$gte": earlyDate,
				"$lte": lastDate,
			},
			"payload": U.Dict{
				"$containsi": "Absent",
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return ctx.Error(500, "Internal Server Error")
	}

	if status != 200 {
		fmt.Println("absentEmployee :", absentEmployee)
		fmt.Println(status)
		return ctx.Error(500, "Internal Server Error")
	}

	return 200, U.Dict{"employee": employee, "checkinEmployee": checkinEmployee, "absentEmployee": absentEmployee}, nil
}
