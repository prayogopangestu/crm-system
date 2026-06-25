package postgres

import (
	"context"
	"strings"
	"time"

	"github.com/prayogopangestu/crm-system/backend/internal/domain"
	"gorm.io/gorm"
)

func (r *Repository) ListContacts(ctx context.Context, organizationID, search, status string, page domain.Page) (domain.ContactList, error) {
	query := r.query(ctx).Model(&contactModel{}).
		Where("organization_id = ? AND deleted_at IS NULL", organizationID)
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ? OR company ILIKE ?", pattern, pattern, pattern)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return domain.ContactList{}, err
	}
	var models []contactModel
	if err := query.
		Order("created_at DESC").
		Limit(page.Limit).
		Offset((page.Page - 1) * page.Limit).
		Find(&models).Error; err != nil {
		return domain.ContactList{}, err
	}
	contacts := make([]domain.Contact, 0, len(models))
	for _, model := range models {
		contact := contactFromModel(model)
		r.decorateContact(&contact)
		contacts = append(contacts, contact)
	}
	return domain.ContactList{Data: contacts, Total: total, Page: page.Page}, nil
}

func (r *Repository) CreateContact(ctx context.Context, principal domain.Principal, input domain.ContactInput) (domain.Contact, error) {
	var contact domain.Contact
	err := r.query(ctx).Transaction(func(tx *gorm.DB) error {
		ownerID := principal.UserID
		model := contactModel{
			OrganizationID:  principal.OrganizationID,
			OwnerID:         &ownerID,
			Name:            input.Name,
			Email:           strings.ToLower(input.Email),
			Company:         input.Company,
			Role:            input.Role,
			Status:          input.Status,
			AvatarURL:       input.AvatarURL,
			LastContactedAt: time.Now(),
		}
		if err := tx.Create(&model).Error; err != nil {
			return mapError(err)
		}
		if err := tx.Create(&activityModel{
			OrganizationID: principal.OrganizationID,
			ActorID:        principal.UserID,
			ActorName:      principal.Name,
			Action:         "menambahkan kontak baru",
			Target:         model.Name,
		}).Error; err != nil {
			return mapError(err)
		}
		contact = contactFromModel(model)
		return nil
	})
	if err != nil {
		return domain.Contact{}, err
	}
	r.decorateContact(&contact)
	return contact, nil
}

func (r *Repository) UpdateContact(ctx context.Context, principal domain.Principal, id string, input domain.ContactInput) (domain.Contact, error) {
	now := time.Now()
	result := r.query(ctx).Model(&contactModel{}).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		Updates(map[string]any{
			"name":              input.Name,
			"email":             strings.ToLower(input.Email),
			"company":           input.Company,
			"role":              input.Role,
			"status":            input.Status,
			"avatar_url":        input.AvatarURL,
			"last_contacted_at": now,
			"updated_at":        now,
		})
	if result.Error != nil {
		return domain.Contact{}, mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.Contact{}, domain.ErrNotFound
	}
	var model contactModel
	if err := r.query(ctx).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		First(&model).Error; err != nil {
		return domain.Contact{}, mapError(err)
	}
	contact := contactFromModel(model)
	r.decorateContact(&contact)
	return contact, nil
}

func (r *Repository) DeleteContact(ctx context.Context, principal domain.Principal, id string) error {
	now := time.Now()
	result := r.query(ctx).Model(&contactModel{}).
		Where("id = ? AND organization_id = ? AND deleted_at IS NULL", id, principal.OrganizationID).
		Updates(map[string]any{"deleted_at": now, "updated_at": now})
	if result.Error != nil {
		return mapError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) decorateContact(contact *domain.Contact) {
	now := time.Now().In(r.location)
	contact.LastContacted = humanTime(contact.LastContactedAt, now)
	contact.Initials = initials(contact.Name)
}
