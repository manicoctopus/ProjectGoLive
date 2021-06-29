package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ProjectGoLive/pkg/models"

	"github.com/gorilla/mux"
)

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		app.errorLog.Print("ERROR occurred while decoding user data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid User data", Description: "Unable to decipher User data"}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	app.infoLog.Printf("adding user name :: %s with email as :: %s"+
		" with hashedPassword as :: %s with contact as :: %s"+
		" with isBOwner as :: %d with isVerified as :: %d "+
		" and created as :: %s", user.UserName,
		user.UserEmail, user.HashedPassword, user.UserContact,
		user.IsBOwner, user.IsVerified, user.Created)

	if id, err := app.users.Create(&user); err != nil {
		app.errorLog.Println(err)
		errMsg := &models.ErrorMsg{Name: "User not created", Description: "Unable to create new User"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		user.UserID = uint32(id)
		json.NewEncoder(w).Encode(user)
	}
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		app.errorLog.Print("ERROR occurred while decoding user data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid User ID", Description: "Unable to decipher User ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.users.Retrieve(user.UserID); err != nil {
			if err != models.ErrNoRecord {
				app.errorLog.Print("ERROR occurred while retriving User data from database :: ", err)
				errMsg := &models.ErrorMsg{Name: "User not found", Description: "User is not in database"}
				json.NewEncoder(w).Encode(errMsg)
			} else {
				app.infoLog.Printf("upserting user name :: %s with email as :: %s"+
					" with hashedPassword as :: %s with contact as :: %s"+
					" with isBOwner as :: %d with isVerified as :: %d "+
					" and created as :: %s", user.UserName,
					user.UserEmail, user.HashedPassword, user.UserContact,
					user.IsBOwner, user.IsVerified, user.Created)

				if _, err := app.users.Create(&user); err != nil {
					app.errorLog.Println(err)
					errMsg := &models.ErrorMsg{Name: "User not upserted", Description: "Unable to upsert User"}
					json.NewEncoder(w).Encode(errMsg)
					return
				}
			}
		} else {
			app.infoLog.Printf("updating user id :: %d with name :: %s with email as :: %s"+
				" with hashedPassword as :: %s with contact as :: %s"+
				" with isBOwner as :: %d with isVerified as :: %d "+
				" and created as :: %s", user.UserID, user.UserName,
				user.UserEmail, user.HashedPassword, user.UserContact,
				user.IsBOwner, user.IsVerified, user.Created)

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
	user := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		app.errorLog.Print("ERROR occurred while decoding user data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid User ID", Description: "Unable to decipher User ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.users.Retrieve(user.UserID); err != nil {
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
			app.infoLog.Printf("deleting user id :: %d with name :: %s with email as :: %s"+
				" with hashedPassword as :: %s with contact as :: %s"+
				" with isBOwner as :: %d with isVerified as :: %d "+
				" and created as :: %s", user.UserID, user.UserName,
				user.UserEmail, user.HashedPassword, user.UserContact,
				user.IsBOwner, user.IsVerified, user.Created)

			if err = app.users.Delete(user.UserID); err != nil {
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
	vars := mux.Vars(r)
	if id, err := strconv.Atoi(vars["id"]); err != nil {
		app.errorLog.Print("ERROR occurred while decoding user data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid User ID", Description: "Unable to decipher User ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if user, err := app.users.Retrieve(uint32(id)); err != nil {
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
	if users, err := app.users.RetrieveAll(); err != nil {
		app.errorLog.Print("ERROR getting users :: ", err)
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
	pdtsvc := models.Pdtsvc{}
	if err := json.NewDecoder(r.Body).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("ERROR occurred while decoding pdtsvc data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Pdtsvc data", Description: "Unable to decipher Pdtsvc data"}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	app.infoLog.Printf("adding pdtsvc name :: %s with price as :: %.2f"+
		" with description as :: %s with catID as :: %d with listID as :: %d"+
		" with views as :: %d with likes as :: %d with keyword as :: %s"+
		" with created as :: %s and modified as :: %s",
		pdtsvc.PdtsvcName, pdtsvc.PdtsvcPrice, pdtsvc.PdtsvcDesc,
		pdtsvc.CatID, pdtsvc.ListID, pdtsvc.Views,
		pdtsvc.Likes, pdtsvc.Keyword, pdtsvc.Created, pdtsvc.Modified)

	if id, err := app.pdtsvcs.Create(&pdtsvc); err != nil {
		app.errorLog.Println(err)
		errMsg := &models.ErrorMsg{Name: "Pdtsvc not created", Description: "Unable to create new Pdtsvc"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		pdtsvc.PdtsvcID = uint32(id)
		json.NewEncoder(w).Encode(pdtsvc)
	}
}

func (app *application) updatePdtsvc(w http.ResponseWriter, r *http.Request) {
	pdtsvc := models.Pdtsvc{}
	if err := json.NewDecoder(r.Body).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("ERROR occurred while decoding pdtsvc data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Pdtsvc ID", Description: "Unable to decipher Pdtsvc ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.pdtsvcs.Retrieve(pdtsvc.PdtsvcID); err != nil {
			if err != models.ErrNoRecord {
				app.errorLog.Print("ERROR occurred while retriving pdtsvc data from database :: ", err)
				errMsg := &models.ErrorMsg{Name: "Pdtsvc not found", Description: "Pdtsvc is not in database"}
				json.NewEncoder(w).Encode(errMsg)
			} else {
				app.infoLog.Printf("upserting pdtsvc name :: %s with price as :: %.2f"+
					" with description as :: %s with catID as :: %d with listID as :: %d"+
					" with views as :: %d with likes as :: %d with keyword as :: %s"+
					" with created as :: %s and modified as :: %s",
					pdtsvc.PdtsvcName, pdtsvc.PdtsvcPrice, pdtsvc.PdtsvcDesc,
					pdtsvc.CatID, pdtsvc.ListID, pdtsvc.Views,
					pdtsvc.Likes, pdtsvc.Keyword, pdtsvc.Created, pdtsvc.Modified)

				if _, err := app.pdtsvcs.Create(&pdtsvc); err != nil {
					app.errorLog.Println(err)
					errMsg := &models.ErrorMsg{Name: "Pdtsvc not upserted", Description: "Unable to upsert Pdtsvc"}
					json.NewEncoder(w).Encode(errMsg)
					return
				}
			}
		} else {
			app.infoLog.Printf("updating pdtsvc id :: %d with name :: %s with price as :: %.2f"+
				" with description as :: %s with catID as :: %d with listID as :: %d"+
				" with views as :: %d with likes as :: %d with keyword as :: %s"+
				" with created as :: %s and modified as :: %s",
				pdtsvc.PdtsvcID, pdtsvc.PdtsvcName, pdtsvc.PdtsvcPrice,
				pdtsvc.PdtsvcDesc, pdtsvc.CatID, pdtsvc.ListID, pdtsvc.Views,
				pdtsvc.Likes, pdtsvc.Keyword, pdtsvc.Created, pdtsvc.Modified)

			if err := app.pdtsvcs.Update(&pdtsvc); err != nil {
				app.errorLog.Println(err)
				errMsg := &models.ErrorMsg{Name: "Pdtsvc not updated", Description: "Unable to update Pdtsvc"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				if tempPdtsvc, err := app.pdtsvcs.Retrieve(pdtsvc.PdtsvcID); err != nil {
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
	pdtsvc := models.Pdtsvc{}
	if err := json.NewDecoder(r.Body).Decode(&pdtsvc); err != nil {
		app.errorLog.Print("ERROR occurred while decoding pdtsvc data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Pdtsvc ID", Description: "Unable to decipher Pdtsvc ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.pdtsvcs.Retrieve(pdtsvc.PdtsvcID); err != nil {
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
			app.infoLog.Printf("deleting pdtsvc id :: %d with name :: %s with price as :: %.2f"+
				" with description as :: %s with catID as :: %d with listID as :: %d"+
				" with views as :: %d with likes as :: %d with keyword as :: %s"+
				" with created as :: %s and modified as :: %s",
				pdtsvc.PdtsvcID, pdtsvc.PdtsvcName, pdtsvc.PdtsvcPrice,
				pdtsvc.PdtsvcDesc, pdtsvc.CatID, pdtsvc.ListID, pdtsvc.Views,
				pdtsvc.Likes, pdtsvc.Keyword, pdtsvc.Created, pdtsvc.Modified)

			if err = app.pdtsvcs.Delete(pdtsvc.PdtsvcID); err != nil {
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
	vars := mux.Vars(r)
	if id, err := strconv.Atoi(vars["id"]); err != nil {
		app.errorLog.Print("ERROR occurred while decoding pdtsvc data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Pdtsvc ID", Description: "Unable to decipher Pdtsvc ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if pdtsvc, err := app.pdtsvcs.Retrieve(uint32(id)); err != nil {
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
	if pdtsvcs, err := app.pdtsvcs.RetrieveAll(); err != nil {
		app.errorLog.Print("ERROR getting pdtsvcs :: ", err)
		errMsg := &models.ErrorMsg{Name: "Pdtsvcs not found", Description: "Unable to retrieve Pdtsvcs from database"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		json.NewEncoder(w).Encode(pdtsvcs)
	}
}

func (app *application) createListing(w http.ResponseWriter, r *http.Request) {
	listing := models.Listing{}
	if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
		app.errorLog.Print("ERROR occurred while decoding listing data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Listing data", Description: "Unable to decipher Listing data"}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	app.infoLog.Printf("adding listing name :: %s with description as :: %s"+
		" with ig_url as :: %s with fb_url as :: %s with website_url as :: %s"+
		" with userID as :: %d with created as :: %s and modified as :: %s",
		listing.ListName, listing.ListDesc, listing.Ig_url,
		listing.Fb_url, listing.Website_url, listing.UserID,
		listing.Created, listing.Modified)

	if id, err := app.listings.Create(&listing); err != nil {
		app.errorLog.Println(err)
		errMsg := &models.ErrorMsg{Name: "Listing not created", Description: "Unable to create new Listing"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		listing.ListID = uint32(id)
		json.NewEncoder(w).Encode(listing)
	}
}

func (app *application) updateListing(w http.ResponseWriter, r *http.Request) {
	listing := models.Listing{}
	if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
		app.errorLog.Print("ERROR occurred while decoding listing data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Listing ID", Description: "Unable to decipher Listing ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.listings.Retrieve(listing.ListID); err != nil {
			if err != models.ErrNoRecord {
				app.errorLog.Print("ERROR occurred while retriving listing data from database :: ", err)
				errMsg := &models.ErrorMsg{Name: "Listing not found", Description: "Listing is not in database"}
				json.NewEncoder(w).Encode(errMsg)
			} else {
				app.infoLog.Printf("upserting listing name :: %s with description as :: %s"+
					" with ig_url as :: %s with fb_url as :: %s with website_url as :: %s"+
					" with userID as :: %d with created as :: %s and modified as :: %s",
					listing.ListName, listing.ListDesc, listing.Ig_url,
					listing.Fb_url, listing.Website_url, listing.UserID,
					listing.Created, listing.Modified)

				if _, err := app.listings.Create(&listing); err != nil {
					app.errorLog.Println(err)
					errMsg := &models.ErrorMsg{Name: "Listing not upserted", Description: "Unable to upsert Listing"}
					json.NewEncoder(w).Encode(errMsg)
					return
				}
			}
		} else {
			app.infoLog.Printf("updating listing id :: %d with name :: %s with description as :: %s"+
				" with ig_url as :: %s with fb_url as :: %s with website_url as :: %s"+
				" with userID as :: %d with created as :: %s and modified as :: %s",
				listing.ListID, listing.ListName, listing.ListDesc, listing.Ig_url,
				listing.Fb_url, listing.Website_url, listing.UserID,
				listing.Created, listing.Modified)

			if err := app.listings.Update(&listing); err != nil {
				app.errorLog.Println(err)
				errMsg := &models.ErrorMsg{Name: "Listing not updated", Description: "Unable to update Listing"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				if tempListing, err := app.listings.Retrieve(listing.ListID); err != nil {
					app.errorLog.Print("ERROR occurred while retriving listing data from database :: ", err)
					errMsg := &models.ErrorMsg{Name: "Listing not found", Description: "Listing is not in database"}
					json.NewEncoder(w).Encode(errMsg)
				} else {
					json.NewEncoder(w).Encode(tempListing)
				}
			}
		}
	}
}

func (app *application) deleteListing(w http.ResponseWriter, r *http.Request) {
	listing := models.Listing{}
	if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
		app.errorLog.Print("ERROR occurred while decoding listing data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Listing ID", Description: "Unable to decipher Listing ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.listings.Retrieve(listing.ListID); err != nil {
			if err == models.ErrNoRecord {
				app.errorLog.Print("ERROR listing data not found in database :: ", err)
				errMsg := &models.ErrorMsg{Name: "Listing not found", Description: "Listing is not in database"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				app.errorLog.Print("ERROR occurred while retriving listing data from database:: ", err)
				errMsg := &models.ErrorMsg{Name: "Listing not found", Description: "Unexpected error retrieving Listing from database"}
				json.NewEncoder(w).Encode(errMsg)
				return
			}
		} else {
			app.infoLog.Printf("deleting listing id :: %d with name :: %s with description as :: %s"+
				" with ig_url as :: %s with fb_url as :: %s with website_url as :: %s"+
				" with userID as :: %d with created as :: %s and modified as :: %s",
				listing.ListID, listing.ListName, listing.ListDesc, listing.Ig_url,
				listing.Fb_url, listing.Website_url, listing.UserID,
				listing.Created, listing.Modified)

			if err = app.listings.Delete(listing.ListID); err != nil {
				app.errorLog.Print("ERROR occurred while deleting listing data from database:: ", err)
				errMsg := &models.ErrorMsg{Name: "Listing not deleted", Description: "Unable to delete Listing"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				app.retrieveAllListings(w, r)
			}
		}
	}
}

func (app *application) retrieveListing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if id, err := strconv.Atoi(vars["id"]); err != nil {
		app.errorLog.Print("ERROR occurred while decoding listing data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Listing ID", Description: "Unable to decipher Listing ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if listing, err := app.listings.Retrieve(uint32(id)); err != nil {
			app.errorLog.Print("ERROR occurred while retriving listing data from database:: ", err)
			errMsg := &models.ErrorMsg{Name: "Listing not found", Description: "Listing is not in database"}
			json.NewEncoder(w).Encode(errMsg)
			return
		} else {
			json.NewEncoder(w).Encode(listing)
		}
	}
}

func (app *application) retrieveAllListings(w http.ResponseWriter, r *http.Request) {
	if listings, err := app.listings.RetrieveAll(); err != nil {
		app.errorLog.Print("ERROR getting listings :: ", err)
		errMsg := &models.ErrorMsg{Name: "Listings not found", Description: "Unable to retrieve Listings from database"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		json.NewEncoder(w).Encode(listings)
	}
}

func (app *application) createReview(w http.ResponseWriter, r *http.Request) {
	review := models.Review{}
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		app.errorLog.Print("ERROR occurred while decoding review data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Review data", Description: "Unable to decipher Review data"}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	app.infoLog.Printf("adding review text :: %s "+
		" with userID as :: %d with listID as :: %d and created as :: %s",
		review.ReviewText, review.UserID, review.ListID, review.Created)

	if id, err := app.reviews.Create(&review); err != nil {
		app.errorLog.Println(err)
		errMsg := &models.ErrorMsg{Name: "Review not created", Description: "Unable to create new Review"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		review.ReviewID = uint32(id)
		json.NewEncoder(w).Encode(review)
	}
}

func (app *application) updateReview(w http.ResponseWriter, r *http.Request) {
	review := models.Review{}
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		app.errorLog.Print("ERROR occurred while decoding review data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Review ID", Description: "Unable to decipher Review ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.reviews.Retrieve(review.ListID); err != nil {
			if err != models.ErrNoRecord {
				app.errorLog.Print("ERROR occurred while retriving review data from database :: ", err)
				errMsg := &models.ErrorMsg{Name: "Review not found", Description: "Review is not in database"}
				json.NewEncoder(w).Encode(errMsg)
			} else {
				app.infoLog.Printf("upserting review text :: %s "+
					" with userID as :: %d with listID as :: %d and created as :: %s",
					review.ReviewText, review.UserID, review.ListID, review.Created)

				if _, err := app.reviews.Create(&review); err != nil {
					app.errorLog.Println(err)
					errMsg := &models.ErrorMsg{Name: "Review not upserted", Description: "Unable to upsert Review"}
					json.NewEncoder(w).Encode(errMsg)
					return
				}
			}
		} else {
			app.infoLog.Printf("updating review id :: %d with text :: %s "+
				" with userID as :: %d with listID as :: %d and created as :: %s",
				review.ReviewID, review.ReviewText, review.UserID, review.ListID,
				review.Created)

			if err := app.reviews.Update(&review); err != nil {
				app.errorLog.Println(err)
				errMsg := &models.ErrorMsg{Name: "Review not updated", Description: "Unable to update Review"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				if tempReview, err := app.reviews.Retrieve(review.ReviewID); err != nil {
					app.errorLog.Print("ERROR occurred while retriving review data from database :: ", err)
					errMsg := &models.ErrorMsg{Name: "Review not found", Description: "Review is not in database"}
					json.NewEncoder(w).Encode(errMsg)
				} else {
					json.NewEncoder(w).Encode(tempReview)
				}
			}
		}
	}
}

func (app *application) deleteReview(w http.ResponseWriter, r *http.Request) {
	review := models.Review{}
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		app.errorLog.Print("ERROR occurred while decoding review data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Review ID", Description: "Unable to decipher Review ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.reviews.Retrieve(review.ReviewID); err != nil {
			if err == models.ErrNoRecord {
				app.errorLog.Print("ERROR review data not found in database :: ", err)
				errMsg := &models.ErrorMsg{Name: "Review not found", Description: "Review is not in database"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				app.errorLog.Print("ERROR occurred while retriving review data from database:: ", err)
				errMsg := &models.ErrorMsg{Name: "Review not found", Description: "Unexpected error retrieving Review from database"}
				json.NewEncoder(w).Encode(errMsg)
				return
			}
		} else {
			app.infoLog.Printf("deleting review id :: %d with text :: %s "+
				" with userID as :: %d with listID as :: %d and created as :: %s",
				review.ReviewID, review.ReviewText, review.UserID, review.ListID,
				review.Created)

			if err = app.reviews.Delete(review.ReviewID); err != nil {
				app.errorLog.Print("ERROR occurred while deleting review data from database:: ", err)
				errMsg := &models.ErrorMsg{Name: "Review not deleted", Description: "Unable to delete Review"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				app.retrieveAllReviews(w, r)
			}
		}
	}
}

func (app *application) retrieveReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if id, err := strconv.Atoi(vars["id"]); err != nil {
		app.errorLog.Print("ERROR occurred while decoding review data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Review ID", Description: "Unable to decipher Review ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if review, err := app.reviews.Retrieve(uint32(id)); err != nil {
			app.errorLog.Print("ERROR occurred while retriving review data from database:: ", err)
			errMsg := &models.ErrorMsg{Name: "Review not found", Description: "Review is not in database"}
			json.NewEncoder(w).Encode(errMsg)
			return
		} else {
			json.NewEncoder(w).Encode(review)
		}
	}
}

func (app *application) retrieveAllReviews(w http.ResponseWriter, r *http.Request) {
	if reviews, err := app.reviews.RetrieveAll(); err != nil {
		app.errorLog.Print("ERROR getting reviews :: ", err)
		errMsg := &models.ErrorMsg{Name: "Review not found", Description: "Unable to retrieve Reviews from database"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		json.NewEncoder(w).Encode(reviews)
	}
}

func (app *application) createCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		app.errorLog.Print("ERROR occurred while decoding category data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Category data", Description: "Unable to decipher Category data"}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	app.infoLog.Printf("adding category name :: %s "+
		" with parentCat as :: %d",
		category.CatName, category.ParentCat)

	if id, err := app.categories.Create(&category); err != nil {
		app.errorLog.Println(err)
		errMsg := &models.ErrorMsg{Name: "Category not created", Description: "Unable to create new Category"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		category.CatID = uint32(id)
		json.NewEncoder(w).Encode(category)
	}
}

func (app *application) updateCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		app.errorLog.Print("ERROR occurred while decoding category data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Category ID", Description: "Unable to decipher Category ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.categories.Retrieve(category.CatID); err != nil {
			if err != models.ErrNoRecord {
				app.errorLog.Print("ERROR occurred while retriving category data from database :: ", err)
				errMsg := &models.ErrorMsg{Name: "Category not found", Description: "Category is not in database"}
				json.NewEncoder(w).Encode(errMsg)
			} else {
				app.infoLog.Printf("upserting category name :: %s "+
					" with parentCat as :: %d",
					category.CatName, category.ParentCat)

				if _, err := app.categories.Create(&category); err != nil {
					app.errorLog.Println(err)
					errMsg := &models.ErrorMsg{Name: "Category not upserted", Description: "Unable to upsert Category"}
					json.NewEncoder(w).Encode(errMsg)
					return
				}
			}
		} else {
			app.infoLog.Printf("updating category id :: %d with name :: %s "+
				" with parentCat as :: %d", category.CatID, category.CatName, category.ParentCat)

			if err := app.categories.Update(&category); err != nil {
				app.errorLog.Println(err)
				errMsg := &models.ErrorMsg{Name: "Category not updated", Description: "Unable to update Category"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				if tempCategory, err := app.categories.Retrieve(category.CatID); err != nil {
					app.errorLog.Print("ERROR occurred while retriving category data from database :: ", err)
					errMsg := &models.ErrorMsg{Name: "Category not found", Description: "Category is not in database"}
					json.NewEncoder(w).Encode(errMsg)
				} else {
					json.NewEncoder(w).Encode(tempCategory)
				}
			}
		}
	}
}

func (app *application) deleteCategory(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		app.errorLog.Print("ERROR occurred while decoding category data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Category ID", Description: "Unable to decipher Category ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if _, err := app.categories.Retrieve(category.CatID); err != nil {
			if err == models.ErrNoRecord {
				app.errorLog.Print("ERROR category data not found in database :: ", err)
				errMsg := &models.ErrorMsg{Name: "Category not found", Description: "Category is not in database"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				app.errorLog.Print("ERROR occurred while retriving category data from database:: ", err)
				errMsg := &models.ErrorMsg{Name: "Category not found", Description: "Unexpected error retrieving Category from database"}
				json.NewEncoder(w).Encode(errMsg)
				return
			}
		} else {
			app.infoLog.Printf("deleting category id :: %d with name :: %s "+
				" with parentCat as :: %d", category.CatID, category.CatName, category.ParentCat)

			if err = app.categories.Delete(category.CatID); err != nil {
				app.errorLog.Print("ERROR occurred while deleting category data from database:: ", err)
				errMsg := &models.ErrorMsg{Name: "Category not deleted", Description: "Unable to delete Category"}
				json.NewEncoder(w).Encode(errMsg)
				return
			} else {
				app.retrieveAllCategories(w, r)
			}
		}
	}
}

func (app *application) retrieveCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if id, err := strconv.Atoi(vars["id"]); err != nil {
		app.errorLog.Print("ERROR occurred while decoding category data :: ", err)
		errMsg := &models.ErrorMsg{Name: "Invalid Category ID", Description: "Unable to decipher Category ID"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		if category, err := app.categories.Retrieve(uint32(id)); err != nil {
			app.errorLog.Print("ERROR occurred while retriving category data from database:: ", err)
			errMsg := &models.ErrorMsg{Name: "Category not found", Description: "Category is not in database"}
			json.NewEncoder(w).Encode(errMsg)
			return
		} else {
			json.NewEncoder(w).Encode(category)
		}
	}
}

func (app *application) retrieveAllCategories(w http.ResponseWriter, r *http.Request) {
	if categories, err := app.categories.RetrieveAll(); err != nil {
		app.errorLog.Print("ERROR getting categories :: ", err)
		errMsg := &models.ErrorMsg{Name: "Category not found", Description: "Unable to retrieve Categories from database"}
		json.NewEncoder(w).Encode(errMsg)
		return
	} else {
		json.NewEncoder(w).Encode(categories)
	}
}
