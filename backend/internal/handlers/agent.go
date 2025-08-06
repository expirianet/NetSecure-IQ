package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/yourorg/yourapp/internal/models"
    "github.com/yourorg/yourapp/internal/wireguard"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

func PreRegisterAgent(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var req struct {
            MacAddress string `json:"mac_address"`
        }
        if err := c.BodyParser(&req); err != nil {
            return fiber.ErrBadRequest
        }
        var agent models.Agent
        if err := db.Where("mac_address = ?", req.MacAddress).First(&agent).Error; err == nil {
            return fiber.NewError(fiber.StatusConflict, "Agent already registered")
        }
        peer, script, err := wireguard.CreatePeerAndScript(req.MacAddress)
        if err != nil {
            return fiber.ErrInternalServerError
        }
        agent = models.Agent{
            ID:              uuid.New(),
            MacAddress:      req.MacAddress,
            Status:          "unassociated",
            WireGuardPeerID: peer.ID,
        }
        if err := db.Create(&agent).Error; err != nil {
            return fiber.ErrInternalServerError
        }
        return c.JSON(fiber.Map{
            "script": script,
            "agent":  agent,
        })
    }
}
