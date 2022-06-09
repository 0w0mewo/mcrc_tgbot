// Code generated by SQLBoiler 4.9.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Chat is an object representing the database table.
type Chat struct {
	ID   int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name string `boil:"name" json:"name" toml:"name" yaml:"name"`

	R *chatR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L chatL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ChatColumns = struct {
	ID   string
	Name string
}{
	ID:   "id",
	Name: "name",
}

var ChatTableColumns = struct {
	ID   string
	Name string
}{
	ID:   "chats.id",
	Name: "chats.name",
}

// Generated where

var ChatWhere = struct {
	ID   whereHelperint64
	Name whereHelperstring
}{
	ID:   whereHelperint64{field: "\"chats\".\"id\""},
	Name: whereHelperstring{field: "\"chats\".\"name\""},
}

// ChatRels is where relationship names are stored.
var ChatRels = struct {
	ChatTweetSubscriptions string
	Messages               string
}{
	ChatTweetSubscriptions: "ChatTweetSubscriptions",
	Messages:               "Messages",
}

// chatR is where relationships are stored.
type chatR struct {
	ChatTweetSubscriptions ChatTweetSubscriptionSlice `boil:"ChatTweetSubscriptions" json:"ChatTweetSubscriptions" toml:"ChatTweetSubscriptions" yaml:"ChatTweetSubscriptions"`
	Messages               MessageSlice               `boil:"Messages" json:"Messages" toml:"Messages" yaml:"Messages"`
}

// NewStruct creates a new relationship struct
func (*chatR) NewStruct() *chatR {
	return &chatR{}
}

// chatL is where Load methods for each relationship are stored.
type chatL struct{}

var (
	chatAllColumns            = []string{"id", "name"}
	chatColumnsWithoutDefault = []string{}
	chatColumnsWithDefault    = []string{"id", "name"}
	chatPrimaryKeyColumns     = []string{"id"}
	chatGeneratedColumns      = []string{"id", "name"}
)

type (
	// ChatSlice is an alias for a slice of pointers to Chat.
	// This should almost always be used instead of []Chat.
	ChatSlice []*Chat
	// ChatHook is the signature for custom Chat hook methods
	ChatHook func(context.Context, boil.ContextExecutor, *Chat) error

	chatQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	chatType                 = reflect.TypeOf(&Chat{})
	chatMapping              = queries.MakeStructMapping(chatType)
	chatPrimaryKeyMapping, _ = queries.BindMapping(chatType, chatMapping, chatPrimaryKeyColumns)
	chatInsertCacheMut       sync.RWMutex
	chatInsertCache          = make(map[string]insertCache)
	chatUpdateCacheMut       sync.RWMutex
	chatUpdateCache          = make(map[string]updateCache)
	chatUpsertCacheMut       sync.RWMutex
	chatUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var chatAfterSelectHooks []ChatHook

var chatBeforeInsertHooks []ChatHook
var chatAfterInsertHooks []ChatHook

var chatBeforeUpdateHooks []ChatHook
var chatAfterUpdateHooks []ChatHook

var chatBeforeDeleteHooks []ChatHook
var chatAfterDeleteHooks []ChatHook

var chatBeforeUpsertHooks []ChatHook
var chatAfterUpsertHooks []ChatHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Chat) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range chatAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Chat) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range chatBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Chat) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range chatAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Chat) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range chatBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Chat) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range chatAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Chat) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range chatBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Chat) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range chatAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Chat) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range chatBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Chat) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range chatAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddChatHook registers your hook function for all future operations.
func AddChatHook(hookPoint boil.HookPoint, chatHook ChatHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		chatAfterSelectHooks = append(chatAfterSelectHooks, chatHook)
	case boil.BeforeInsertHook:
		chatBeforeInsertHooks = append(chatBeforeInsertHooks, chatHook)
	case boil.AfterInsertHook:
		chatAfterInsertHooks = append(chatAfterInsertHooks, chatHook)
	case boil.BeforeUpdateHook:
		chatBeforeUpdateHooks = append(chatBeforeUpdateHooks, chatHook)
	case boil.AfterUpdateHook:
		chatAfterUpdateHooks = append(chatAfterUpdateHooks, chatHook)
	case boil.BeforeDeleteHook:
		chatBeforeDeleteHooks = append(chatBeforeDeleteHooks, chatHook)
	case boil.AfterDeleteHook:
		chatAfterDeleteHooks = append(chatAfterDeleteHooks, chatHook)
	case boil.BeforeUpsertHook:
		chatBeforeUpsertHooks = append(chatBeforeUpsertHooks, chatHook)
	case boil.AfterUpsertHook:
		chatAfterUpsertHooks = append(chatAfterUpsertHooks, chatHook)
	}
}

// One returns a single chat record from the query.
func (q chatQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Chat, error) {
	o := &Chat{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for chats")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Chat records from the query.
func (q chatQuery) All(ctx context.Context, exec boil.ContextExecutor) (ChatSlice, error) {
	var o []*Chat

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Chat slice")
	}

	if len(chatAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Chat records in the query.
func (q chatQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count chats rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q chatQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if chats exists")
	}

	return count > 0, nil
}

// ChatTweetSubscriptions retrieves all the chat_tweet_subscription's ChatTweetSubscriptions with an executor.
func (o *Chat) ChatTweetSubscriptions(mods ...qm.QueryMod) chatTweetSubscriptionQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"chat_tweet_subscriptions\".\"chat_id\"=?", o.ID),
	)

	return ChatTweetSubscriptions(queryMods...)
}

// Messages retrieves all the message's Messages with an executor.
func (o *Chat) Messages(mods ...qm.QueryMod) messageQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"messages\".\"chat_id\"=?", o.ID),
	)

	return Messages(queryMods...)
}

// LoadChatTweetSubscriptions allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (chatL) LoadChatTweetSubscriptions(ctx context.Context, e boil.ContextExecutor, singular bool, maybeChat interface{}, mods queries.Applicator) error {
	var slice []*Chat
	var object *Chat

	if singular {
		object = maybeChat.(*Chat)
	} else {
		slice = *maybeChat.(*[]*Chat)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &chatR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &chatR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`chat_tweet_subscriptions`),
		qm.WhereIn(`chat_tweet_subscriptions.chat_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load chat_tweet_subscriptions")
	}

	var resultSlice []*ChatTweetSubscription
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice chat_tweet_subscriptions")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on chat_tweet_subscriptions")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for chat_tweet_subscriptions")
	}

	if len(chatTweetSubscriptionAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.ChatTweetSubscriptions = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &chatTweetSubscriptionR{}
			}
			foreign.R.Chat = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.ChatID) {
				local.R.ChatTweetSubscriptions = append(local.R.ChatTweetSubscriptions, foreign)
				if foreign.R == nil {
					foreign.R = &chatTweetSubscriptionR{}
				}
				foreign.R.Chat = local
				break
			}
		}
	}

	return nil
}

// LoadMessages allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (chatL) LoadMessages(ctx context.Context, e boil.ContextExecutor, singular bool, maybeChat interface{}, mods queries.Applicator) error {
	var slice []*Chat
	var object *Chat

	if singular {
		object = maybeChat.(*Chat)
	} else {
		slice = *maybeChat.(*[]*Chat)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &chatR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &chatR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`messages`),
		qm.WhereIn(`messages.chat_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load messages")
	}

	var resultSlice []*Message
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice messages")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on messages")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for messages")
	}

	if len(messageAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Messages = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &messageR{}
			}
			foreign.R.Chat = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.ChatID) {
				local.R.Messages = append(local.R.Messages, foreign)
				if foreign.R == nil {
					foreign.R = &messageR{}
				}
				foreign.R.Chat = local
				break
			}
		}
	}

	return nil
}

// AddChatTweetSubscriptions adds the given related objects to the existing relationships
// of the chat, optionally inserting them as new records.
// Appends related to o.R.ChatTweetSubscriptions.
// Sets related.R.Chat appropriately.
func (o *Chat) AddChatTweetSubscriptions(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*ChatTweetSubscription) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.ChatID, o.ID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"chat_tweet_subscriptions\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 0, []string{"chat_id"}),
				strmangle.WhereClause("\"", "\"", 0, chatTweetSubscriptionPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.ChatID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &chatR{
			ChatTweetSubscriptions: related,
		}
	} else {
		o.R.ChatTweetSubscriptions = append(o.R.ChatTweetSubscriptions, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &chatTweetSubscriptionR{
				Chat: o,
			}
		} else {
			rel.R.Chat = o
		}
	}
	return nil
}

// SetChatTweetSubscriptions removes all previously related items of the
// chat replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Chat's ChatTweetSubscriptions accordingly.
// Replaces o.R.ChatTweetSubscriptions with related.
// Sets related.R.Chat's ChatTweetSubscriptions accordingly.
func (o *Chat) SetChatTweetSubscriptions(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*ChatTweetSubscription) error {
	query := "update \"chat_tweet_subscriptions\" set \"chat_id\" = null where \"chat_id\" = ?"
	values := []interface{}{o.ID}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, values)
	}
	_, err := exec.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.ChatTweetSubscriptions {
			queries.SetScanner(&rel.ChatID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Chat = nil
		}
		o.R.ChatTweetSubscriptions = nil
	}

	return o.AddChatTweetSubscriptions(ctx, exec, insert, related...)
}

// RemoveChatTweetSubscriptions relationships from objects passed in.
// Removes related items from R.ChatTweetSubscriptions (uses pointer comparison, removal does not keep order)
// Sets related.R.Chat.
func (o *Chat) RemoveChatTweetSubscriptions(ctx context.Context, exec boil.ContextExecutor, related ...*ChatTweetSubscription) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.ChatID, nil)
		if rel.R != nil {
			rel.R.Chat = nil
		}
		if err = rel.Update(ctx, exec, boil.Whitelist("chat_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.ChatTweetSubscriptions {
			if rel != ri {
				continue
			}

			ln := len(o.R.ChatTweetSubscriptions)
			if ln > 1 && i < ln-1 {
				o.R.ChatTweetSubscriptions[i] = o.R.ChatTweetSubscriptions[ln-1]
			}
			o.R.ChatTweetSubscriptions = o.R.ChatTweetSubscriptions[:ln-1]
			break
		}
	}

	return nil
}

// AddMessages adds the given related objects to the existing relationships
// of the chat, optionally inserting them as new records.
// Appends related to o.R.Messages.
// Sets related.R.Chat appropriately.
func (o *Chat) AddMessages(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Message) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.ChatID, o.ID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"messages\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 0, []string{"chat_id"}),
				strmangle.WhereClause("\"", "\"", 0, messagePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			queries.Assign(&rel.ChatID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &chatR{
			Messages: related,
		}
	} else {
		o.R.Messages = append(o.R.Messages, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &messageR{
				Chat: o,
			}
		} else {
			rel.R.Chat = o
		}
	}
	return nil
}

// SetMessages removes all previously related items of the
// chat replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Chat's Messages accordingly.
// Replaces o.R.Messages with related.
// Sets related.R.Chat's Messages accordingly.
func (o *Chat) SetMessages(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Message) error {
	query := "update \"messages\" set \"chat_id\" = null where \"chat_id\" = ?"
	values := []interface{}{o.ID}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, values)
	}
	_, err := exec.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.Messages {
			queries.SetScanner(&rel.ChatID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Chat = nil
		}
		o.R.Messages = nil
	}

	return o.AddMessages(ctx, exec, insert, related...)
}

// RemoveMessages relationships from objects passed in.
// Removes related items from R.Messages (uses pointer comparison, removal does not keep order)
// Sets related.R.Chat.
func (o *Chat) RemoveMessages(ctx context.Context, exec boil.ContextExecutor, related ...*Message) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.ChatID, nil)
		if rel.R != nil {
			rel.R.Chat = nil
		}
		if err = rel.Update(ctx, exec, boil.Whitelist("chat_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Messages {
			if rel != ri {
				continue
			}

			ln := len(o.R.Messages)
			if ln > 1 && i < ln-1 {
				o.R.Messages[i] = o.R.Messages[ln-1]
			}
			o.R.Messages = o.R.Messages[:ln-1]
			break
		}
	}

	return nil
}

// Chats retrieves all the records using an executor.
func Chats(mods ...qm.QueryMod) chatQuery {
	mods = append(mods, qm.From("\"chats\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"chats\".*"})
	}

	return chatQuery{NewQuery(mods...)}
}

// FindChat retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindChat(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*Chat, error) {
	chatObj := &Chat{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"chats\" where \"id\"=?", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, chatObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from chats")
	}

	if err = chatObj.doAfterSelectHooks(ctx, exec); err != nil {
		return chatObj, err
	}

	return chatObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Chat) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no chats provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(chatColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	chatInsertCacheMut.RLock()
	cache, cached := chatInsertCache[key]
	chatInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			chatAllColumns,
			chatColumnsWithDefault,
			chatColumnsWithoutDefault,
			nzDefaults,
		)
		wl = strmangle.SetComplement(wl, chatGeneratedColumns)

		cache.valueMapping, err = queries.BindMapping(chatType, chatMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(chatType, chatMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"chats\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"chats\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into chats")
	}

	if !cached {
		chatInsertCacheMut.Lock()
		chatInsertCache[key] = cache
		chatInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Chat.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Chat) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return err
	}
	key := makeCacheKey(columns, nil)
	chatUpdateCacheMut.RLock()
	cache, cached := chatUpdateCache[key]
	chatUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			chatAllColumns,
			chatPrimaryKeyColumns,
		)
		wl = strmangle.SetComplement(wl, chatGeneratedColumns)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update chats, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"chats\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 0, wl),
			strmangle.WhereClause("\"", "\"", 0, chatPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(chatType, chatMapping, append(wl, chatPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	_, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update chats row")
	}

	if !cached {
		chatUpdateCacheMut.Lock()
		chatUpdateCache[key] = cache
		chatUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q chatQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for chats")
	}

	return nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ChatSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), chatPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"chats\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, chatPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in chat slice")
	}

	return nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Chat) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no chats provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(chatColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	chatUpsertCacheMut.RLock()
	cache, cached := chatUpsertCache[key]
	chatUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			chatAllColumns,
			chatColumnsWithDefault,
			chatColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			chatAllColumns,
			chatPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert chats, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(chatPrimaryKeyColumns))
			copy(conflict, chatPrimaryKeyColumns)
		}
		cache.query = buildUpsertQuerySQLite(dialect, "\"chats\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(chatType, chatMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(chatType, chatMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert chats")
	}

	if !cached {
		chatUpsertCacheMut.Lock()
		chatUpsertCache[key] = cache
		chatUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Chat record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Chat) Delete(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil {
		return errors.New("models: no Chat provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), chatPrimaryKeyMapping)
	sql := "DELETE FROM \"chats\" WHERE \"id\"=?"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from chats")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return err
	}

	return nil
}

// DeleteAll deletes all matching rows.
func (q chatQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if q.Query == nil {
		return errors.New("models: no chatQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from chats")
	}

	return nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ChatSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) error {
	if len(o) == 0 {
		return nil
	}

	if len(chatBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), chatPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"chats\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, chatPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	_, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from chat slice")
	}

	if len(chatAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Chat) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindChat(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ChatSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ChatSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), chatPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"chats\".* FROM \"chats\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, chatPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ChatSlice")
	}

	*o = slice

	return nil
}

// ChatExists checks if the Chat row exists.
func ChatExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"chats\" where \"id\"=? limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if chats exists")
	}

	return exists, nil
}