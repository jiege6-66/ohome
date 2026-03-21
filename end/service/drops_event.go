package service

import (
	"errors"
	"fmt"
	"ohome/global"
	"ohome/model"
	"ohome/service/dto"
	"sort"
	"strings"
	"time"
)

type DropsEventService struct {
	BaseService
}

func (s *DropsEventService) GetList(listDTO *dto.DropsEventListDTO, loginUser model.LoginUser, summaryOnly bool, loc *time.Location) ([]model.DropsEvent, int64, error) {
	query := global.DB.Model(&model.DropsEvent{}).
		Where("scope_type = ? OR (scope_type = ? AND owner_user_id = ?)", model.DropsScopeShared, model.DropsScopePersonal, loginUser.ID)

	scopeType := strings.TrimSpace(strings.ToLower(listDTO.ScopeType))
	switch scopeType {
	case model.DropsScopeShared:
		query = query.Where("scope_type = ?", model.DropsScopeShared)
	case model.DropsScopePersonal:
		query = query.Where("scope_type = ? AND owner_user_id = ?", model.DropsScopePersonal, loginUser.ID)
	}
	if eventType := normalizeDropsEventType(listDTO.EventType); eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if keyword := strings.TrimSpace(listDTO.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("(title LIKE ? OR remark LIKE ?)", like, like)
	}

	events := make([]model.DropsEvent, 0, 32)
	if err := query.Order("updated_at DESC").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	now := time.Now().In(loc)
	filtered := make([]model.DropsEvent, 0, len(events))
	for _, event := range events {
		if summaryOnly && !event.Enabled {
			continue
		}
		nextOccurAt, err := computeDropsEventNextOccurrence(&event, now, loc)
		if err != nil {
			continue
		}
		event.NextOccurAt = nextOccurAt
		if listDTO.Month > 0 {
			if nextOccurAt == nil || int(nextOccurAt.Month()) != listDTO.Month {
				continue
			}
		}
		filtered = append(filtered, event)
	}

	sort.Slice(filtered, func(i, j int) bool {
		left := filtered[i].NextOccurAt
		right := filtered[j].NextOccurAt
		switch {
		case left == nil && right == nil:
			return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt)
		case left == nil:
			return false
		case right == nil:
			return true
		default:
			return left.Before(*right)
		}
	})

	total := int64(len(filtered))
	start := (listDTO.GetPage() - 1) * listDTO.GetLimit()
	if start >= len(filtered) {
		return []model.DropsEvent{}, total, nil
	}
	end := start + listDTO.GetLimit()
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (s *DropsEventService) GetByID(id uint, loginUser model.LoginUser, loc *time.Location) (model.DropsEvent, error) {
	var event model.DropsEvent
	if err := global.DB.First(&event, id).Error; err != nil {
		return model.DropsEvent{}, err
	}
	if event.ScopeType != model.DropsScopeShared && event.OwnerUserID != loginUser.ID {
		return model.DropsEvent{}, errors.New("无权限访问该重要日期")
	}
	nextOccurAt, _ := computeDropsEventNextOccurrence(&event, time.Now().In(loc), loc)
	event.NextOccurAt = nextOccurAt
	return event, nil
}

func (s *DropsEventService) AddOrUpdate(updateDTO *dto.DropsEventUpsertDTO, loginUser model.LoginUser, loc *time.Location) (model.DropsEvent, error) {
	scopeType := normalizeDropsScope(updateDTO.ScopeType)
	eventType := normalizeDropsEventType(updateDTO.EventType)
	if eventType == "" {
		return model.DropsEvent{}, errors.New("重要日期类型无效")
	}
	calendarType := normalizeDropsCalendarType(updateDTO.CalendarType)
	reminderRaw, _, err := configuredDropsEventReminderDays()
	if err != nil {
		return model.DropsEvent{}, err
	}

	repeatYearly := true
	if updateDTO.RepeatYearly != nil {
		repeatYearly = *updateDTO.RepeatYearly
	}
	enabled := true
	if updateDTO.Enabled != nil {
		enabled = *updateDTO.Enabled
	}

	event := model.DropsEvent{}
	if updateDTO.ID != 0 {
		existing, err := s.GetByID(updateDTO.ID, loginUser, loc)
		if err != nil {
			return model.DropsEvent{}, err
		}
		event = existing
	}

	updateDTO.ApplyToModel(&event)
	event.ScopeType = scopeType
	event.Title = strings.TrimSpace(event.Title)
	event.EventType = eventType
	event.CalendarType = calendarType
	event.OwnerUserID = loginUser.ID
	if event.CreatedBy == 0 {
		event.CreatedBy = loginUser.ID
	}
	event.UpdatedBy = loginUser.ID
	event.Remark = strings.TrimSpace(event.Remark)
	event.ReminderDays = reminderRaw
	event.RepeatYearly = repeatYearly
	event.Enabled = enabled

	if event.Title == "" {
		return model.DropsEvent{}, errors.New("重要日期标题不能为空")
	}
	if err := validateDropsEvent(&event, loc); err != nil {
		return model.DropsEvent{}, err
	}
	if err := global.DB.Save(&event).Error; err != nil {
		return model.DropsEvent{}, err
	}
	savedEvent, err := s.GetByID(event.ID, loginUser, loc)
	if err != nil {
		return model.DropsEvent{}, err
	}
	TriggerDropsReminderRescan(loc)
	return savedEvent, nil
}

func (s *DropsEventService) Delete(id uint, loginUser model.LoginUser, loc *time.Location) error {
	event, err := s.GetByID(id, loginUser, loc)
	if err != nil {
		return err
	}
	return global.DB.Delete(&model.DropsEvent{}, event.ID).Error
}

func (s *DropsEventService) ListReminderCandidates(loc *time.Location) ([]model.DropsEvent, error) {
	events := make([]model.DropsEvent, 0, 32)
	if err := global.DB.Model(&model.DropsEvent{}).Where("enabled = ?", true).Find(&events).Error; err != nil {
		return nil, err
	}
	now := time.Now().In(loc)
	result := make([]model.DropsEvent, 0, len(events))
	for _, event := range events {
		nextOccurAt, err := computeDropsEventNextOccurrence(&event, now, loc)
		if err != nil || nextOccurAt == nil {
			continue
		}
		event.NextOccurAt = nextOccurAt
		result = append(result, event)
	}
	return result, nil
}

func (s *DropsEventService) BuildReminderMessages(now time.Time, loc *time.Location) ([]model.AppMessage, error) {
	events, err := s.ListReminderCandidates(loc)
	if err != nil {
		return nil, err
	}
	_, reminderDays, err := configuredDropsEventReminderDays()
	if err != nil {
		return nil, err
	}
	userIDs, err := (&AppMessageService{}).ListRecipientUserIDs()
	if err != nil {
		return nil, err
	}
	result := make([]model.AppMessage, 0, 16)
	for _, event := range events {
		if event.NextOccurAt == nil {
			continue
		}
		days := dropsDaysUntil(now, *event.NextOccurAt, loc)
		if !containsInt(reminderDays, days) {
			continue
		}
		recipients := []uint{event.OwnerUserID}
		if event.ScopeType == model.DropsScopeShared {
			recipients = userIDs
		}
		for _, recipientID := range recipients {
			sourceKey := fmt.Sprintf(
				"%s:%d:%s:%d:%d",
				model.DropsBizTypeEvent,
				event.ID,
				event.NextOccurAt.In(loc).Format("2006-01-02"),
				days,
				recipientID,
			)
			result = append(result, model.AppMessage{
				OwnerUserID: recipientID,
				Source:      model.AppMessageSourceDrops,
				SourceKey:   sourceKey,
				MessageType: model.AppMessageTypeDropsEventReminder,
				BizType:     model.DropsBizTypeEvent,
				BizID:       event.ID,
				UniqueKey:   model.BuildAppMessageUniqueKey(model.AppMessageSourceDrops, sourceKey),
				Title:       s.buildEventReminderTitle(event, days),
				Summary:     s.buildEventReminderSummary(event, loc),
				TriggerDate: dropsStartOfDay(*event.NextOccurAt, loc),
			})
		}
	}
	return result, nil
}

func (s *DropsEventService) buildEventReminderTitle(event model.DropsEvent, days int) string {
	if days == 0 {
		return fmt.Sprintf("%s 就是今天", event.Title)
	}
	return fmt.Sprintf("%s 将在 %d 天后到来", event.Title, days)
}

func (s *DropsEventService) buildEventReminderSummary(event model.DropsEvent, loc *time.Location) string {
	parts := []string{"类型：" + dropsEventTypeLabel(event.EventType)}
	if event.NextOccurAt != nil {
		parts = append(parts, "日期："+event.NextOccurAt.In(loc).Format("2006-01-02"))
	}
	return strings.Join(parts, "；")
}

func validateDropsEvent(event *model.DropsEvent, loc *time.Location) error {
	if event.EventMonth < 1 || event.EventMonth > 12 {
		return errors.New("月份无效")
	}
	if event.EventDay < 1 || event.EventDay > 31 {
		return errors.New("日期无效")
	}
	if !event.RepeatYearly && event.EventYear <= 0 {
		return errors.New("非循环重要日期必须填写年份")
	}
	if _, err := computeDropsEventNextOccurrence(event, time.Now().In(loc), loc); err != nil {
		return err
	}
	return nil
}
