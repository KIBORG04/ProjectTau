package repository

func SetToNullFlavorById(id int32) {
	Database.Exec("update manifest_entries set flavor = null where id = ?;", id)
}
