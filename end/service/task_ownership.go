package service

func EnsureTaskOwnership() error {
	ownerUserID, err := userDao.GetPreferredTaskOwnerID()
	if err != nil || ownerUserID == 0 {
		return err
	}

	if err := quarkAutoSaveTaskDao.AssignMissingOwners(ownerUserID); err != nil {
		return err
	}

	records, err := quarkTransferTaskDao.GetUnownedBySourceTask()
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.SourceTaskID == nil || *record.SourceTaskID == 0 {
			continue
		}
		sourceOwnerUserID, err := quarkAutoSaveTaskDao.GetOwnerUserIDByID(*record.SourceTaskID)
		if err != nil || sourceOwnerUserID == 0 {
			continue
		}
		if err := quarkTransferTaskDao.UpdateOwnerByID(record.ID, sourceOwnerUserID); err != nil {
			return err
		}
	}

	return quarkTransferTaskDao.AssignMissingOwners(ownerUserID)
}
