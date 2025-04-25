#!/bin/bash

# Kullanım: ./generate-structure.sh ModelAdi "Field1 type1" "Field2 type2" ...

# Yardım mesajı
function print_help {
    echo "Kullanım: $0 ModelAdi 'Field1 type1' 'Field2 type2' ..."
    echo ""
    echo "ModelAdi: Oluşturulacak modelin adı"
    echo "Field1, Field2: Modelin alanları ve türleri"
    echo "Örnek:"
    echo "  ./generate-structure.sh User 'Name string' 'Age int'"
    echo "  Bu komut, User modeli için Name ve Age alanlarıyla dosya yapısını oluşturur."
}

# Eğer --help veya -h parametresi verilirse, yardım mesajı göster
if [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]] || [ "$#" -lt 2 ]; then
    print_help
    exit 0
fi

MODEL_NAME=$1
FIELDS=("${@:2}")  # Alanlar
LOWER_MODEL_NAME=$(echo "$MODEL_NAME" | tr '[:upper:]' '[:lower:]')  # Küçük harfe çevir
BASE_DIR="internal"

if [ -f "${BASE_DIR}/model/${LOWER_MODEL_NAME}.go" ]; then
    echo "Uyarı: '$MODEL_NAME' adlı model zaten mevcut. İşlem başlatılmayacak."
    exit 1
fi

# Model dosyası yoksa, işlemi başlat
echo "Model '$MODEL_NAME' oluşturuluyor..."

# Model dosyasını oluştur
cat <<EOL > $BASE_DIR/model/${LOWER_MODEL_NAME}.go
package model

import (
    "github.com/uptrace/bun"
    "time"
)

type $MODEL_NAME struct {
    bun.BaseModel \`bun:"table:${LOWER_MODEL_NAME}s,alias:${LOWER_MODEL_NAME}"\`

    ID        int64     \`json:"id" bun:",pk,autoincrement"\`
    CreatedAt time.Time \`json:"created_at" bun:",nullzero,default:current_timestamp"\`
    UpdatedAt time.Time \`json:"updated_at" bun:",nullzero,default:current_timestamp"\`
EOL

# Alanları ekle
for FIELD in "${FIELDS[@]}"; do
    FIELD_NAME=$(echo "$FIELD" | cut -d' ' -f1)
    FIELD_TYPE=$(echo "$FIELD" | cut -d' ' -f2)
    cat <<EOL >> $BASE_DIR/model/${LOWER_MODEL_NAME}.go
    ${FIELD_NAME} ${FIELD_TYPE} \`json:"${FIELD_NAME}"\`
EOL
done

cat <<EOL >> $BASE_DIR/model/${LOWER_MODEL_NAME}.go
}
EOL

# DTO dosyası oluştur
cat <<EOL > $BASE_DIR/dto/${LOWER_MODEL_NAME}_dto.go
package dto

import (
  "time"
)

type Create${MODEL_NAME}Request struct {
EOL

for FIELD in "${FIELDS[@]}"; do
    FIELD_NAME=$(echo "$FIELD" | cut -d' ' -f1)
    FIELD_TYPE=$(echo "$FIELD" | cut -d' ' -f2)
    cat <<EOL >> $BASE_DIR/dto/${LOWER_MODEL_NAME}_dto.go
    ${FIELD_NAME} ${FIELD_TYPE} \`json:"${FIELD_NAME}" validate:"required"\`
EOL
done

cat <<EOL >> $BASE_DIR/dto/${LOWER_MODEL_NAME}_dto.go
}

func (dto Create${MODEL_NAME}Request) ToDBModel(m model.${MODEL_NAME}) model.${MODEL_NAME} {
EOL

for FIELD in "${FIELDS[@]}"; do
    FIELD_NAME=$(echo "$FIELD" | cut -d' ' -f1)
    cat <<EOL >> $BASE_DIR/dto/${LOWER_MODEL_NAME}_dto.go
    m.${FIELD_NAME} = dto.${FIELD_NAME}
EOL
done

cat <<EOL >> $BASE_DIR/dto/${LOWER_MODEL_NAME}_dto.go
    return m
}

type Update${MODEL_NAME}Request struct {
EOL

for FIELD in "${FIELDS[@]}"; do
    FIELD_NAME=$(echo "$FIELD" | cut -d' ' -f1)
    FIELD_TYPE=$(echo "$FIELD" | cut -d' ' -f2)
    cat <<EOL >> internal/dto/${LOWER_MODEL_NAME}_dto.go
    ${FIELD_NAME} ${FIELD_TYPE} \`json:"${FIELD_NAME}"\`
EOL
done

cat <<EOL >> internal/dto/${LOWER_MODEL_NAME}_dto.go
}

func (dto Update${MODEL_NAME}Request) ToDBModel(m model.${MODEL_NAME}) model.${MODEL_NAME} {
EOL

for FIELD in "${FIELDS[@]}"; do
    FIELD_NAME=$(echo "$FIELD" | cut -d' ' -f1)
    cat <<EOL >> $BASE_DIR/dto/${LOWER_MODEL_NAME}_dto.go
    m.${FIELD_NAME} = dto.${FIELD_NAME}
EOL
done

cat <<EOL >> $BASE_DIR/dto/${LOWER_MODEL_NAME}_dto.go
    return m
}

type ${MODEL_NAME}Response struct {
    ID        int64     \`json:"id"\`
    CreatedAt time.Time \`json:"created_at"\`
    UpdatedAt time.Time \`json:"updated_at"\`
EOL

for FIELD in "${FIELDS[@]}"; do
    FIELD_NAME=$(echo "$FIELD" | cut -d' ' -f1)
    FIELD_TYPE=$(echo "$FIELD" | cut -d' ' -f2)
    cat <<EOL >> internal/dto/${LOWER_MODEL_NAME}_dto.go
    ${FIELD_NAME} ${FIELD_TYPE} \`json:"${FIELD_NAME}"\`
EOL
done

cat <<EOL >> internal/dto/${LOWER_MODEL_NAME}_dto.go
}

func (dto ${MODEL_NAME}Response) ToResponseModel(m model.${MODEL_NAME}) ${MODEL_NAME}Response {
    dto.ID = m.ID
    dto.CreatedAt = m.CreatedAt
    dto.UpdatedAt = m.UpdatedAt
EOL

for FIELD in "${FIELDS[@]}"; do
    FIELD_NAME=$(echo "$FIELD" | cut -d' ' -f1)
    cat <<EOL >> $BASE_DIR/dto/${LOWER_MODEL_NAME}_dto.go
    dto.${FIELD_NAME} = m.${FIELD_NAME}
EOL
done

cat <<EOL >> $BASE_DIR/dto/${LOWER_MODEL_NAME}_dto.go
    return dto
}
EOL

# Repository dosyası oluştur
cat <<EOL > $BASE_DIR/repository/${LOWER_MODEL_NAME}_repository.go
package repository

import (
    "context"
    "github.com/uptrace/bun"
)

type I${MODEL_NAME}Repository interface {
    Create(ctx context.Context, ${LOWER_MODEL_NAME} *model.${MODEL_NAME}) error
    GetByID(ctx context.Context, id int64) (*model.${MODEL_NAME}, error)
    Update(ctx context.Context, ${LOWER_MODEL_NAME} *model.${MODEL_NAME}) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context) ([]model.${MODEL_NAME}, error)
}

type ${MODEL_NAME}Repository struct {
    db *bun.DB
}

func New${MODEL_NAME}Repository(db *bun.DB) I${MODEL_NAME}Repository {
    return &${MODEL_NAME}Repository{db: db}
}

func (r *${MODEL_NAME}Repository) Create(ctx context.Context, ${LOWER_MODEL_NAME} *model.${MODEL_NAME}) error {
    _, err := r.db.NewInsert().Model(${LOWER_MODEL_NAME}).Exec(ctx)
    return err
}

func (r *${MODEL_NAME}Repository) GetByID(ctx context.Context, id int64) (*model.${MODEL_NAME}, error) {
    var ${LOWER_MODEL_NAME} model.${MODEL_NAME}
    err := r.db.NewSelect().Model(&${LOWER_MODEL_NAME}).Where("id = ?", id).Scan(ctx)
    return &${LOWER_MODEL_NAME}, err
}

func (r *${MODEL_NAME}Repository) Update(ctx context.Context, ${LOWER_MODEL_NAME} *model.${MODEL_NAME}) error {
    _, err := r.db.NewUpdate().Model(${LOWER_MODEL_NAME}).WherePK().Exec(ctx)
    return err
}

func (r *${MODEL_NAME}Repository) Delete(ctx context.Context, id int64) error {
    _, err := r.db.NewDelete().Model((*model.${MODEL_NAME})(nil)).Where("id = ?", id).Exec(ctx)
    return err
}

func (r *${MODEL_NAME}Repository) List(ctx context.Context) ([]model.${MODEL_NAME}, error) {
    var ${LOWER_MODEL_NAME}s []model.${MODEL_NAME}
    err := r.db.NewSelect().Model(&${LOWER_MODEL_NAME}s).Scan(ctx)
    return ${LOWER_MODEL_NAME}s, err
}
EOL

# Service dosyası oluştur
cat <<EOL > $BASE_DIR/service/${LOWER_MODEL_NAME}_service.go
package service

import (
    "context"
)

type ${MODEL_NAME}Service struct {
    ${LOWER_MODEL_NAME}Repo repository.I${MODEL_NAME}Repository
}

func New${MODEL_NAME}Service(repo repository.I${MODEL_NAME}Repository) *${MODEL_NAME}Service {
    return &${MODEL_NAME}Service{${LOWER_MODEL_NAME}Repo: repo}
}

func (s *${MODEL_NAME}Service) Create(ctx context.Context, ${LOWER_MODEL_NAME} *model.${MODEL_NAME}) error {
    if err := s.${LOWER_MODEL_NAME}repo.Create(ctx, ${LOWER_MODEL_NAME}); err != nil {
        return errorx.Wrap(errorx.ErrDatabaseOperation, err)
    }
    return nil
}

func (s *${MODEL_NAME}Service) GetByID(ctx context.Context, id int64) (*model.${MODEL_NAME}, error) {
    ${LOWER_MODEL_NAME}, err := s.${LOWER_MODEL_NAME}Repo.GetByID(ctx, id)
    if err != nil {
        return nil, errorx.Wrap(errorx.ErrDatabaseOperation, err)
    }
    return ${LOWER_MODEL_NAME}, nil
}

func (s *${MODEL_NAME}Service) Update(ctx context.Context, ${LOWER_MODEL_NAME} model.${MODEL_NAME}) error {
    if err := s.${LOWER_MODEL_NAME}Repo.Update(ctx, &${LOWER_MODEL_NAME}); err != nil {
        return errorx.Wrap(errorx.ErrDatabaseOperation, err)
    }
    return nil
}

func (s *${MODEL_NAME}Service) Delete(ctx context.Context, id int64) error {
    if err := s.${LOWER_MODEL_NAME}Repo.Delete(ctx, id); err != nil {
        return errorx.Wrap(errorx.ErrDatabaseOperation, err)
    }
    return nil
}

func (s *${MODEL_NAME}Service) List(ctx context.Context) ([]model.${MODEL_NAME}, error) {
    ${LOWER_MODEL_NAME}s, err := s.${LOWER_MODEL_NAME}Repo.List(ctx)
    if err != nil {
        return nil, errorx.Wrap(errorx.ErrDatabaseOperation, err)
    }
    return ${LOWER_MODEL_NAME}s, nil
}
EOL

# Handler dosyası oluştur
cat <<EOL > $BASE_DIR/handler/${LOWER_MODEL_NAME}_handler.go
package handler

import (
  "github.com/gofiber/fiber/v2"
)

type ${MODEL_NAME}Handler struct {
    service *service.${MODEL_NAME}Service
}

func New${MODEL_NAME}Handler(s *service.${MODEL_NAME}Service) *${MODEL_NAME}Handler {
    return &${MODEL_NAME}Handler{service: s}
}

func (h *${MODEL_NAME}Handler) Create(c *fiber.Ctx) error {
    var req dto.Create${MODEL_NAME}Request
    if err := c.BodyParser(&req); err != nil {
        return errorx.Wrap(errorx.ErrInvalidRequest, err)
    }

    ${LOWER_MODEL_NAME} := req.ToDBModel(model.${MODEL_NAME}{})

    if err := h.service.Create(c.Context(),&${LOWER_MODEL_NAME}); err != nil {
        return errorx.WithDetails(errorx.ErrInternal, err.Error())
    }

    return response.Success(c, nil, "$MODEL_NAME başarıyla oluşturuldu")
}

func (h *${MODEL_NAME}Handler) GetByID(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return errorx.Wrap(errorx.ErrInvalidRequest, err)
    }

    resp, err := h.service.GetByID(c.Context(), int64(id))
    if err != nil {
      return errorx.WithDetails(errorx.ErrNotFound, "$MODEL_NAME bulunamadı")
    }

    ${LOWER_MODEL_NAME} := dto.${MODEL_NAME}Response{}.ToResponseModel(*resp)

	  return response.Success(c, ${LOWER_MODEL_NAME})
}

func (h *${MODEL_NAME}Handler) Update(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return errorx.Wrap(errorx.ErrInvalidRequest, err)
    }

    var req dto.Update${MODEL_NAME}Request
    if err = c.BodyParser(&req); err != nil {
        return errorx.Wrap(errorx.ErrInvalidRequest, err)
    }

    _, err = h.service.GetByID(c.Context(), int64(id))
    if err != nil {
        return err
    }

    ${LOWER_MODEL_NAME} := req.ToDBModel(model.${MODEL_NAME}{})

    if err = h.service.Update(c.Context(), ${LOWER_MODEL_NAME}); err != nil {
        return errorx.WithDetails(errorx.ErrInternal, err.Error())
    }

	return response.Success(c, nil, "$MODEL_NAME başarıyla güncellendi")
}

func (h *${MODEL_NAME}Handler) Delete(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return errorx.Wrap(errorx.ErrInvalidRequest, err)
    }

    if err = h.service.Delete(c.Context(), int64(id)); err != nil {
        return errorx.WithDetails(errorx.ErrInternal, err.Error())
    }

	return response.Success(c, nil, "$MODEL_NAME başarıyla silindi")
}

func (h *${MODEL_NAME}Handler) List(c *fiber.Ctx) error {
    resp, err := h.service.List(c.Context())
    if err != nil {
        return errorx.WithDetails(errorx.ErrInternal, err.Error())
    }

    ${LOWER_MODEL_NAME}s := make([]dto.${MODEL_NAME}Response, len(resp))
    for i, item := range resp {
        ${LOWER_MODEL_NAME}s[i] = dto.${MODEL_NAME}Response{}.ToResponseModel(item)
    }
	  return response.Success(c, ${LOWER_MODEL_NAME}s)
}
EOL

echo "$MODEL_NAME için dosyalar oluşturuldu!"
