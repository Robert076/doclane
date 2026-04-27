package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Robert076/doclane/backend/models"
	"github.com/Robert076/doclane/backend/repositories"
	"github.com/Robert076/doclane/backend/types"
	errors "github.com/Robert076/doclane/backend/types/errors"
)

const defaultTagColor = "#6366f1"

type TagService struct {
	tagRepo repositories.ITagRepo
	logger  *slog.Logger
}

func NewTagService(tagRepo repositories.ITagRepo, logger *slog.Logger) *TagService {
	return &TagService{
		tagRepo: tagRepo,
		logger:  logger,
	}
}

func (s *TagService) GetTags(ctx context.Context) ([]models.Tag, error) {
	tags, err := s.tagRepo.GetTags(ctx)
	if err != nil {
		s.logger.Error("failed to get tags", slog.Any("error", err))
		return nil, err
	}
	return tags, nil
}

func (s *TagService) GetTagByID(ctx context.Context, id int) (models.Tag, error) {
	tag, err := s.tagRepo.GetTagByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get tag",
			slog.Int("tag_id", id),
			slog.Any("error", err),
		)
		return models.Tag{}, err
	}
	return tag, nil
}

func (s *TagService) CreateTag(ctx context.Context, claims types.JWTClaims, dto models.TagDTOCreate) (models.Tag, error) {
	if !claims.IsAdmin() {
		s.logger.Warn("non-admin attempted to create tag",
			slog.Int("jwt_user_id", claims.UserID),
		)
		return models.Tag{}, errors.ErrForbidden{Msg: "Only admins can manage tags."}
	}

	tags, err := s.tagRepo.GetTags(ctx)
	if err != nil {
		s.logger.Error("error when retrieving tags for checking total count",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return models.Tag{}, err
	}

	MAX_TAGS := 30
	if len(tags) > MAX_TAGS {
		s.logger.Warn(fmt.Sprintf("max tag count of %d has been reached", MAX_TAGS),
			slog.Int("jwt_user_id", claims.UserID),
		)
		return models.Tag{}, errors.ErrBadRequest{Msg: "Maximum tag count has been reached."}
	}

	if err := validateTagDTO(dto.Name, dto.Color); err != nil {
		return models.Tag{}, err
	}

	if dto.Color == "" {
		dto.Color = defaultTagColor
	}

	tag, err := s.tagRepo.CreateTag(ctx, dto)
	if err != nil {
		s.logger.Error("failed to create tag",
			slog.String("name", dto.Name),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return models.Tag{}, err
	}

	s.logger.Info("tag created",
		slog.Int("tag_id", tag.ID),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return tag, nil
}

func (s *TagService) UpdateTag(ctx context.Context, claims types.JWTClaims, id int, dto models.TagDTOUpdate) (models.Tag, error) {
	if !claims.IsAdmin() {
		s.logger.Warn("non-admin attempted to update tag",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("tag_id", id),
		)
		return models.Tag{}, errors.ErrForbidden{Msg: "Only admins can manage tags."}
	}

	if err := validateTagDTO(dto.Name, dto.Color); err != nil {
		return models.Tag{}, err
	}

	tag, err := s.tagRepo.UpdateTag(ctx, id, dto)
	if err != nil {
		s.logger.Error("failed to update tag",
			slog.Int("tag_id", id),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return models.Tag{}, err
	}

	s.logger.Info("tag updated",
		slog.Int("tag_id", id),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return tag, nil
}

func (s *TagService) DeleteTag(ctx context.Context, claims types.JWTClaims, id int) error {
	if !claims.IsAdmin() {
		s.logger.Warn("non-admin attempted to delete tag",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("tag_id", id),
		)
		return errors.ErrForbidden{Msg: "Only admins can manage tags."}
	}

	if err := s.tagRepo.DeleteTag(ctx, id); err != nil {
		s.logger.Error("failed to delete tag",
			slog.Int("tag_id", id),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("tag deleted",
		slog.Int("tag_id", id),
		slog.Int("jwt_user_id", claims.UserID),
	)
	return nil
}

func (s *TagService) SetTemplateTags(ctx context.Context, claims types.JWTClaims, templateID int, tagIDs []int) error {
	if !claims.IsAdmin() {
		s.logger.Warn("non-admin attempted to set template tags",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("template_id", templateID),
		)
		return errors.ErrForbidden{Msg: "Only admins can manage tags."}
	}

	MAX_LENGTH_TAGS := 3
	if len(tagIDs) > MAX_LENGTH_TAGS {
		s.logger.Warn("maximum tag count reached for template",
			slog.Int("jwt_user_id", claims.UserID),
			slog.Int("template_id", templateID),
		)
		return errors.ErrBadRequest{Msg: fmt.Sprintf("Templates can have up to %d tags.", MAX_LENGTH_TAGS)}
	}

	if err := s.tagRepo.SetTemplateTags(ctx, templateID, tagIDs); err != nil {
		s.logger.Error("failed to set template tags",
			slog.Int("template_id", templateID),
			slog.Int("jwt_user_id", claims.UserID),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("template tags updated",
		slog.Int("template_id", templateID),
		slog.Int("jwt_user_id", claims.UserID),
		slog.Int("tag_count", len(tagIDs)),
	)
	return nil
}

func (s *TagService) GetTagsByTemplateID(ctx context.Context, templateID int) ([]models.Tag, error) {
	tags, err := s.tagRepo.GetTagsByTemplateID(ctx, templateID)
	if err != nil {
		s.logger.Error("failed to get tags for template",
			slog.Int("template_id", templateID),
			slog.Any("error", err),
		)
		return nil, err
	}
	return tags, nil
}
