package seed

import (
	"encoding/json"
	"log"
	"time"

	"ukm-hub/internal/entity"
	"ukm-hub/internal/utils"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	var count int64
	if err := db.Model(&entity.Organization{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		log.Println("[seed] data already exists, skipping")
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		superAdmin, err := createUser(tx, "Super Admin", "admin@ukmhub.app", "admin123", utils.RoleSuperAdmin)
		if err != nil {
			return err
		}
		orgAdmin, err := createUser(tx, "Admin CSA", "admin@csa.app", "admin123", utils.RoleOrgAdmin)
		if err != nil {
			return err
		}

		org := entity.Organization{
			Name:        "Creative Student Association",
			Slug:        "creative-student-association",
			Description: "Organisasi mahasiswa kreatif dengan divisi Programming dan Multimedia.",
			Email:       "csa@ukmhub.app",
			Phone:       "081234567890",
			Status:      entity.OrganizationStatusActive,
		}
		if err := tx.Create(&org).Error; err != nil {
			return err
		}

		if err := tx.Create(&entity.OrganizationAdmin{
			UserID:         orgAdmin.ID,
			OrganizationID: org.ID,
		}).Error; err != nil {
			return err
		}

		divisionSpecs := []struct {
			Name        string
			Description string
		}{
			{Name: "Programming", Description: "Divisi pengembangan perangkat lunak dan logika."},
			{Name: "Multimedia", Description: "Divisi desain, video, dan konten kreatif."},
		}
		for _, ds := range divisionSpecs {
			div := entity.Division{
				OrganizationID: org.ID,
				Name:           ds.Name,
				Description:    ds.Description,
				Status:         entity.DivisionStatusActive,
			}
			if err := tx.Create(&div).Error; err != nil {
				return err
			}
		}

		regStart := time.Now().AddDate(0, 0, -7)
		regEnd := time.Now().AddDate(0, 0, 30)
		event := entity.Event{
			OrganizationID:    org.ID,
			Name:              "Open Recruitment CSA 2026",
			Slug:              "open-recruitment-csa-2026",
			Description:       "Pendaftaran anggota baru Creative Student Association tahun 2026.",
			Location:          "Kampus",
			StartDate:         &regStart,
			EndDate:           &regEnd,
			RegistrationStart: &regStart,
			RegistrationEnd:   &regEnd,
			Quota:             100,
			Status:            entity.EventStatusPublished,
		}
		if err := tx.Create(&event).Error; err != nil {
			return err
		}

		form := entity.Form{
			EventID:     event.ID,
			Title:       "Formulir Pendaftaran CSA 2026",
			Description: "Isi formulir berikut untuk mendaftar menjadi anggota CSA.",
			IsPublished: true,
		}
		if err := tx.Create(&form).Error; err != nil {
			return err
		}

		var opts = func(v []string) string {
			b, _ := json.Marshal(v)
			return string(b)
		}

		fields := []entity.FormField{
			{FormID: form.ID, Label: "Nama Lengkap", Name: "nama_lengkap", Type: entity.FormFieldTypeText, Placeholder: "Nama lengkap kamu", Required: true, SortOrder: 1},
			{FormID: form.ID, Label: "NIM", Name: "nim", Type: entity.FormFieldTypeText, Placeholder: "Nomor Induk Mahasiswa", Required: true, SortOrder: 2},
			{FormID: form.ID, Label: "Email", Name: "email", Type: entity.FormFieldTypeEmail, Placeholder: "email@kampus.ac.id", Required: true, SortOrder: 3},
			{FormID: form.ID, Label: "Nomor WhatsApp", Name: "whatsapp", Type: entity.FormFieldTypePhone, Placeholder: "08xxxxxxxxxx", Required: true, SortOrder: 4},
			{FormID: form.ID, Label: "Divisi Pilihan", Name: "divisi_pilihan", Type: entity.FormFieldTypeSelect, Options: opts([]string{"Programming", "Multimedia"}), Required: true, SortOrder: 5},
			{FormID: form.ID, Label: "Alasan Bergabung", Name: "alasan_bergabung", Type: entity.FormFieldTypeTextarea, Required: true, SortOrder: 6},
			{FormID: form.ID, Label: "Portfolio", Name: "portfolio", Type: entity.FormFieldTypeURL, Placeholder: "https://", Required: false, SortOrder: 7},
		}
		for i := range fields {
			if err := tx.Create(&fields[i]).Error; err != nil {
				return err
			}
		}

		log.Println("[seed] seeded demo data (super_admin:", superAdmin.Email, ", org_admin:", orgAdmin.Email, ", org:", org.Slug, ", event:", event.Slug, ")")
		return nil
	})
}

func createUser(tx *gorm.DB, name, email, password, role string) (*entity.User, error) {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := entity.User{
		Name:     name,
		Email:    email,
		Password: hashed,
		Role:     role,
	}
	if err := tx.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
