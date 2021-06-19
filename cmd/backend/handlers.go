package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ProjectGoLive/pkg/models"

	"github.com/gorilla/mux"
)

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Create User\n"))

	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		app.errorLog.Print("ERROR occurred while decoding user data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid User data", Description: "Unable to decipher User data"}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	app.infoLog.Printf("adding user name :: %s with email"+
		" as :: %s and hashedPassword as :: %s", user.Name,
		user.Email, user.HashedPassword)

	if id, err := app.users.Create(&user); err != nil {
		app.errorLog.Println(err)
		errMsg := &models.ErrorMsg{Name: "User not created", Description: "Unable to create new User"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		user.ID = id
		json.NewEncoder(w).Encode(user)
	}
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Update User\n"))

	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		app.errorLog.Print("ERROR occurred while decoding user data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid User ID", Description: "Unable to User Pdtsvc ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.users.Retrieve(user.ID); err != nil {
			if err != models.ErrNoRecord {
				app.errorLog.Print("ERROR occurred while retriving pdtsvc data from database :: ", err)
				errMsg := &models.ErrorMsg{Name: "User not found", Description: "User is not in database"}
				json.NewEncoder(w).Encode(errMsg)
			} else {
				app.infoLog.Printf("upserting user id :: %d with name :: "+
					"%s with email as :: %s and hashedPassword as :: %s ",
					user.ID, user.Name, user.Email, user.HashedPassword)

				if _, err := app.users.Create(&user); err != nil {
					app.errorLog.Println(err)
					errMsg := &models.ErrorMsg{Name: "User not upserted", Description: "Unable to upsert User"}
					json.NewEncoder(w).Encode(errMsg)
					return
				}
			}
		} else {
			app.infoLog.Printf("updating user id :: %d with name :: "+
				"%s with email as :: %s and hashedPassword as :: %s ",
				user.ID, user.Name, user.Email, user.HashedPassword)

			if err := app.users.Update(&user); err != nil {
				app.errorLog.Println(err)
				errMsg := &models.ErrorMsg{Name: "User not updated", Description: "Unable to update User"}
				json.NewEncoder(w).Encode(errMsg)
				return
			}
		}
	}
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Delete User\n"))

	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		app.errorLog.Print("ERROR occurred while decoding user data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid User ID", Description: "Unable to decipher User ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.users.Retrieve(user.ID); err != nil {
			if err == models.ErrNoRecord {
				app.errorLog.Print("ERROR user data not found in database :: ", err)
				errMsg := &models.ErrorMsg{Name: "User not found", Description: "User is not in database"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				app.errorLog.Print("ERROR occurred while retriving user data from database:: ", err)
				errMsg := &models.ErrorMsg{Name: "User not found", Description: "Unexpected error retrieving User from database"}
				json.NewEncoder(w).Encode(errMsg)
				return
			}
		} else {
			app.infoLog.Printf("deleting user id :: %d with user name "+
				":: %s with email as :: %s and hashedPassword as :: %s",
				user.ID, user.Name, user.Email, user.HashedPassword)

			if err = app.users.Delete(user.ID); err != nil {
				app.errorLog.Print("ERROR occurred while deleting user data from database:: ", err)
				errMsg := &models.ErrorMsg{Name: "User not deleted", Description: "Unable to delete User"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				app.retrieveAllUsers(w, r)
			}
		}
	}
}

func (app *application) retrieveUser(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Retrieve User\n"))

	vars := mux.Vars(r)
	if id, err := strconv.Atoi(vars["id"]); err != nil {
		app.errorLog.Print("ERROR occurred while decoding user data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid User ID", Description: "Unable to decipher User ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if user, err := app.users.Retrieve(id); err != nil {
			app.errorLog.Print("ERROR occurred while retriving user data from database:: ", err)
			errMsg := &models.ErrorMsg{Name: "User not found", Description: "User is not in database"}
			json.NewEncoder(w).Encode(errMsg)
			return
		} else {
			json.NewEncoder(w).Encode(user)
		}
	}
}

func (app *application) retrieveAllUsers(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Retrieve Users\n"))

	if users, err := app.users.RetrieveAll(); err != nil {
		app.errorLog.Print("ERROR getting pdtsvcs :: ", err)
		errMsg := &models.ErrorMsg{Name: "Users not found", Description: "Unable to retrieve Users from database"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		json.NewEncoder(w).Encode(users)
	}
}

func (app *application) authenticateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Authenticate User\n"))
}

func (app *application) createPdtsvc(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Create Pdtsvc\n"))

	pdtsvc := models.Pdtsvc{}
	if err := json.NewDecoder(r.Body).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("ERROR occurred while decoding pdtsvc data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Pdtsvc data", Description: "Unable to decipher Pdtsvc data"}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	app.infoLog.Printf("adding pdtsvc code :: %s with name"+
		" as :: %s and description as :: %s", pdtsvc.Code,
		pdtsvc.Name, pdtsvc.Description)

	if id, err := app.pdtsvcs.Create(&pdtsvc); err != nil {
		app.errorLog.Println(err)
		errMsg := &models.ErrorMsg{Name: "Pdtsvc not created", Description: "Unable to create new Pdtsvc"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		pdtsvc.ID = id
		json.NewEncoder(w).Encode(pdtsvc)
	}
}

func (app *application) updatePdtsvc(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Update Pdtsvc\n"))

	pdtsvc := models.Pdtsvc{}
	if err := json.NewDecoder(r.Body).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("ERROR occurred while decoding pdtsvc data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Pdtsvc ID", Description: "Unable to decipher Pdtsvc ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.pdtsvcs.Retrieve(pdtsvc.ID); err != nil {
			if err != models.ErrNoRecord {
				app.errorLog.Print("ERROR occurred while retriving pdtsvc data from database :: ", err)
				errMsg := &models.ErrorMsg{Name: "Pdtsvc not found", Description: "Pdtsvc is not in database"}
				json.NewEncoder(w).Encode(errMsg)
			} else {
				app.infoLog.Printf("upserting pdtsvc id :: %d with"+
					" pdtsvc code :: %s with name"+
					" as :: %s and description as :: %s", pdtsvc.ID,
					pdtsvc.Code, pdtsvc.Name, pdtsvc.Description)

				if _, err := app.pdtsvcs.Create(&pdtsvc); err != nil {
					app.errorLog.Println(err)
					errMsg := &models.ErrorMsg{Name: "Pdtsvc not upserted", Description: "Unable to upsert Pdtsvc"}
					json.NewEncoder(w).Encode(errMsg)
					return
				}
			}
		} else {
			app.infoLog.Printf("updating pdtsvc id :: %d with"+
				" pdtsvc code :: %s with name"+
				" as :: %s and description as :: %s", pdtsvc.ID,
				pdtsvc.Code, pdtsvc.Name, pdtsvc.Description)

			if err := app.pdtsvcs.Update(&pdtsvc); err != nil {
				app.errorLog.Println(err)
				errMsg := &models.ErrorMsg{Name: "Pdtsvc not updated", Description: "Unable to update Pdtsvc"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				if tempPdtsvc, err := app.pdtsvcs.Retrieve(pdtsvc.ID); err != nil {
					app.errorLog.Print("ERROR occurred while retriving pdtsvc data from database :: ", err)
					errMsg := &models.ErrorMsg{Name: "Pdtsvc not found", Description: "Pdtsvc is not in database"}
					json.NewEncoder(w).Encode(errMsg)
				} else {
					json.NewEncoder(w).Encode(tempPdtsvc)
				}
			}
		}
	}
}

func (app *application) deletePdtsvc(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Delete Pdtsvc\n"))

	pdtsvc := models.Pdtsvc{}
	if err := json.NewDecoder(r.Body).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("ERROR occurred while decoding pdtsvc data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Pdtsvc ID", Description: "Unable to decipher Pdtsvc ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.pdtsvcs.Retrieve(pdtsvc.ID); err != nil {
			if err == models.ErrNoRecord {
				app.errorLog.Print("ERROR pdtsvc data not found in database :: ", err)
				errMsg := &models.ErrorMsg{Name: "Pdtsvc not found", Description: "Pdtsvc is not in database"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				app.errorLog.Print("ERROR occurred while retriving pdtsvc data from database:: ", err)
				errMsg := &models.ErrorMsg{Name: "Pdtsvc not found", Description: "Unexpected error retrieving Pdtsvc from database"}
				json.NewEncoder(w).Encode(errMsg)
				return
			}
		} else {
			app.infoLog.Printf("deleting pdtsvc id :: %d with"+
				" pdtsvc code :: %s with name"+
				" as :: %s and description as :: %s", pdtsvc.ID,
				pdtsvc.Code, pdtsvc.Name, pdtsvc.Description)

			if err = app.pdtsvcs.Delete(pdtsvc.ID); err != nil {
				app.errorLog.Print("ERROR occurred while deleting pdtsvc data from database:: ", err)
				errMsg := &models.ErrorMsg{Name: "Pdtsvc not deleted", Description: "Unable to delete Pdtsvc"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				app.retrieveAllPdtsvcs(w, r)
			}
		}
	}
}

func (app *application) retrievePdtsvc(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Retrieve Pdtsvc\n"))

	vars := mux.Vars(r)
	if id, err := strconv.Atoi(vars["id"]); err != nil {
		app.errorLog.Print("ERROR occurred while decoding pdtsvc data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Pdtsvc ID", Description: "Unable to decipher Pdtsvc ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if pdtsvc, err := app.pdtsvcs.Retrieve(id); err != nil {
			app.errorLog.Print("ERROR occurred while retriving pdtsvc data from database:: ", err)
			errMsg := &models.ErrorMsg{Name: "Pdtsvc not found", Description: "Pdtsvc is not in database"}
			json.NewEncoder(w).Encode(errMsg)
			return
		} else {
			json.NewEncoder(w).Encode(pdtsvc)
		}
	}
}

func (app *application) retrieveAllPdtsvcs(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Retrieve Pdtsvcs\n"))

	if pdtsvcs, err := app.pdtsvcs.RetrieveAll(); err != nil {
		app.errorLog.Print("ERROR getting pdtsvcs :: ", err)
		errMsg := &models.ErrorMsg{Name: "Pdtsvcs not found", Description: "Unable to retrieve Pdtsvcs from database"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		json.NewEncoder(w).Encode(pdtsvcs)
	}
}
