// MIT License
// Copyright (c) 2025 Toni Liesche
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

package model

import (
	"cloud-toolbox/internal/infrastructure/errors"
	"github.com/rabbitmq/amqp091-go"
	"strconv"
	"time"
)

type RabbitMQRecord struct {
	AppId           string            `json:"app_id,omitempty"`
	Body            string            `json:"body"`
	ContentEncoding string            `json:"content_encoding,omitempty"`
	ContentType     string            `json:"content_type,omitempty"`
	CorrelationId   string            `json:"correlation_id,omitempty"`
	ConsumerTag     string            `json:"consumer_tag,omitempty"`
	DeliveryMode    uint8             `json:"delivery_mode,omitempty"`
	DeliveryTag     uint64            `json:"delivery_tag,omitempty"`
	Headers         map[string]string `json:"headers,omitempty"`
	Exchange        string            `json:"exchange,omitempty"`
	Expiration      string            `json:"expiration,omitempty"`
	MessageCount    uint32            `json:"message_count,omitempty"`
	MessageId       string            `json:"message_id,omitempty"`
	Priority        uint8             `json:"priority,omitempty"`
	Redelivered     bool              `json:"redelivered"`
	ReplyTo         string            `json:"reply_to,omitempty"`
	RoutingKey      string            `json:"routing_key,omitempty"`
	Timestamp       time.Time         `json:"timestamp,omitempty"`
	Type            string            `json:"type,omitempty"`
	UserId          string            `json:"user_id,omitempty"`
}

func (r *RabbitMQRecord) GetRecordIdentifier() string {
	return r.MessageId
}

func (r *RabbitMQRecord) GetPayload() string {
	return r.Body
}

func (r *RabbitMQRecord) GetMetaInformation() map[string]string {
	meta := make(map[string]string)
	if r.AppId != "" {
		meta["rmq.meta.app_id"] = r.AppId
	}

	if r.ContentEncoding != "" {
		meta["rmq.meta.content_encoding"] = r.ContentEncoding
	}

	if r.ContentType != "" {
		meta["rmq.meta.content_type"] = r.ContentType
	}

	if r.CorrelationId != "" {
		meta["rmq.meta.correlation_id"] = r.CorrelationId
	}

	if r.ConsumerTag != "" {
		meta["rmq.meta.consumer_tag"] = r.ConsumerTag
	}

	if r.DeliveryMode != 0 {
		meta["rmq.meta.delivery_mode"] = strconv.FormatUint(uint64(r.DeliveryMode), 10)
	}

	if r.DeliveryTag != 0 {
		meta["rmq.meta.delivery_tag"] = strconv.FormatUint(r.DeliveryTag, 10)
	}

	if r.Exchange != "" {
		meta["rmq.meta.exchange"] = r.Exchange
	}

	if r.Expiration != "" {
		meta["rmq.meta.expiration"] = r.Expiration
	}

	if r.MessageCount != 0 {
		meta["rmq.meta.message_count"] = strconv.FormatUint(uint64(r.MessageCount), 10)
	}

	if r.Priority != 0 {
		meta["rmq.meta.priority"] = strconv.FormatUint(uint64(r.Priority), 10)
	}

	meta["rmq.meta.redelivered"] = strconv.FormatBool(r.Redelivered)

	if r.ReplyTo != "" {
		meta["rmq.meta.reply_to"] = r.ReplyTo
	}

	if r.RoutingKey != "" {
		meta["rmq.meta.routing_key"] = r.RoutingKey
	}

	nilTime := time.Time{}
	if r.Timestamp != nilTime {
		meta["rmq.meta.timestamp"] = r.Timestamp.String()
	}

	if r.Type != "" {
		meta["rmq.meta.type"] = r.Type
	}

	if r.UserId != "" {
		meta["rmq.meta.user_id"] = r.UserId
	}

	return meta
}

func (r *RabbitMQRecord) Validate() errors.ApplicationError {
	if r.MessageId == "" {
		return errors.NewMissingRequestFieldError("messageId")
	}

	if string(r.Body) == "" {
		return errors.NewMissingRequestFieldError("body")
	}

	return nil
}

func NewRabbitMQRecordFromDelivery(deliver amqp091.Delivery) *RabbitMQRecord {
	headers := make(map[string]string)

	for k, v := range deliver.Headers {
		if str, ok := v.(string); ok {
			headers[k] = str
		} else if bytes, ok := v.([]byte); ok {
			headers[k] = string(bytes)
		}
	}

	return &RabbitMQRecord{
		AppId:           deliver.AppId,
		Body:            string(deliver.Body),
		ContentEncoding: deliver.ContentEncoding,
		ContentType:     deliver.ContentType,
		CorrelationId:   deliver.CorrelationId,
		ConsumerTag:     deliver.ConsumerTag,
		DeliveryMode:    deliver.DeliveryMode,
		DeliveryTag:     deliver.DeliveryTag,
		Exchange:        deliver.Exchange,
		Expiration:      deliver.Expiration,
		Headers:         headers,
		MessageCount:    deliver.MessageCount,
		MessageId:       deliver.MessageId,
		Priority:        deliver.Priority,
		Redelivered:     deliver.Redelivered,
		ReplyTo:         deliver.ReplyTo,
		RoutingKey:      deliver.RoutingKey,
		Timestamp:       deliver.Timestamp,
		Type:            deliver.Type,
		UserId:          deliver.UserId,
	}
}
