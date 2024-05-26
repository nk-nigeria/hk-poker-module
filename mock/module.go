package mock

import (
	"context"
	"os"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/rtapi"
	"github.com/heroiclabs/nakama-common/runtime"
)

var _ runtime.NakamaModule = &MockModule{}

type MockModule struct {
}

// AccountDeleteId implements runtime.NakamaModule.
func (m *MockModule) AccountDeleteId(ctx context.Context, userID string, recorded bool) error {
	panic("unimplemented")
}

// AccountExportId implements runtime.NakamaModule.
func (m *MockModule) AccountExportId(ctx context.Context, userID string) (string, error) {
	panic("unimplemented")
}

// AccountGetId implements runtime.NakamaModule.
func (m *MockModule) AccountGetId(ctx context.Context, userID string) (*api.Account, error) {
	panic("unimplemented")
}

// AccountUpdateId implements runtime.NakamaModule.
func (m *MockModule) AccountUpdateId(ctx context.Context, userID string, username string, metadata map[string]interface{}, displayName string, timezone string, location string, langTag string, avatarUrl string) error {
	panic("unimplemented")
}

// AccountsGetId implements runtime.NakamaModule.
func (m *MockModule) AccountsGetId(ctx context.Context, userIDs []string) ([]*api.Account, error) {
	panic("unimplemented")
}

// AuthenticateApple implements runtime.NakamaModule.
func (m *MockModule) AuthenticateApple(ctx context.Context, token string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateCustom implements runtime.NakamaModule.
func (m *MockModule) AuthenticateCustom(ctx context.Context, id string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateDevice implements runtime.NakamaModule.
func (m *MockModule) AuthenticateDevice(ctx context.Context, id string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateEmail implements runtime.NakamaModule.
func (m *MockModule) AuthenticateEmail(ctx context.Context, email string, password string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateFacebook implements runtime.NakamaModule.
func (m *MockModule) AuthenticateFacebook(ctx context.Context, token string, importFriends bool, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateFacebookInstantGame implements runtime.NakamaModule.
func (m *MockModule) AuthenticateFacebookInstantGame(ctx context.Context, signedPlayerInfo string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateGameCenter implements runtime.NakamaModule.
func (m *MockModule) AuthenticateGameCenter(ctx context.Context, playerID string, bundleID string, timestamp int64, salt string, signature string, publicKeyUrl string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateGoogle implements runtime.NakamaModule.
func (m *MockModule) AuthenticateGoogle(ctx context.Context, token string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateSteam implements runtime.NakamaModule.
func (m *MockModule) AuthenticateSteam(ctx context.Context, token string, username string, create bool) (string, string, bool, error) {
	panic("unimplemented")
}

// AuthenticateTokenGenerate implements runtime.NakamaModule.
func (m *MockModule) AuthenticateTokenGenerate(userID string, username string, exp int64, vars map[string]string) (string, int64, error) {
	panic("unimplemented")
}

// ChannelIdBuild implements runtime.NakamaModule.
func (m *MockModule) ChannelIdBuild(ctx context.Context, sender string, target string, chanType runtime.ChannelType) (string, error) {
	panic("unimplemented")
}

// ChannelMessageRemove implements runtime.NakamaModule.
func (m *MockModule) ChannelMessageRemove(ctx context.Context, channelId string, messageId string, senderId string, senderUsername string, persist bool) (*rtapi.ChannelMessageAck, error) {
	panic("unimplemented")
}

// ChannelMessageSend implements runtime.NakamaModule.
func (m *MockModule) ChannelMessageSend(ctx context.Context, channelID string, content map[string]interface{}, senderId string, senderUsername string, persist bool) (*rtapi.ChannelMessageAck, error) {
	panic("unimplemented")
}

// ChannelMessageUpdate implements runtime.NakamaModule.
func (m *MockModule) ChannelMessageUpdate(ctx context.Context, channelID string, messageID string, content map[string]interface{}, senderId string, senderUsername string, persist bool) (*rtapi.ChannelMessageAck, error) {
	panic("unimplemented")
}

// ChannelMessagesList implements runtime.NakamaModule.
func (m *MockModule) ChannelMessagesList(ctx context.Context, channelId string, limit int, forward bool, cursor string) (messages []*api.ChannelMessage, nextCursor string, prevCursor string, err error) {
	panic("unimplemented")
}

// Event implements runtime.NakamaModule.
func (m *MockModule) Event(ctx context.Context, evt *api.Event) error {
	panic("unimplemented")
}

// FriendsAdd implements runtime.NakamaModule.
func (m *MockModule) FriendsAdd(ctx context.Context, userID string, username string, ids []string, usernames []string) error {
	panic("unimplemented")
}

// FriendsBlock implements runtime.NakamaModule.
func (m *MockModule) FriendsBlock(ctx context.Context, userID string, username string, ids []string, usernames []string) error {
	panic("unimplemented")
}

// FriendsDelete implements runtime.NakamaModule.
func (m *MockModule) FriendsDelete(ctx context.Context, userID string, username string, ids []string, usernames []string) error {
	panic("unimplemented")
}

// FriendsList implements runtime.NakamaModule.
func (m *MockModule) FriendsList(ctx context.Context, userID string, limit int, state *int, cursor string) ([]*api.Friend, string, error) {
	panic("unimplemented")
}

// GetSatori implements runtime.NakamaModule.
func (m *MockModule) GetSatori() runtime.Satori {
	panic("unimplemented")
}

// GroupCreate implements runtime.NakamaModule.
func (m *MockModule) GroupCreate(ctx context.Context, userID string, name string, creatorID string, langTag string, description string, avatarUrl string, open bool, metadata map[string]interface{}, maxCount int) (*api.Group, error) {
	panic("unimplemented")
}

// GroupDelete implements runtime.NakamaModule.
func (m *MockModule) GroupDelete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GroupUpdate implements runtime.NakamaModule.
func (m *MockModule) GroupUpdate(ctx context.Context, id string, userID string, name string, creatorID string, langTag string, description string, avatarUrl string, open bool, metadata map[string]interface{}, maxCount int) error {
	panic("unimplemented")
}

// GroupUserJoin implements runtime.NakamaModule.
func (m *MockModule) GroupUserJoin(ctx context.Context, groupID string, userID string, username string) error {
	panic("unimplemented")
}

// GroupUserLeave implements runtime.NakamaModule.
func (m *MockModule) GroupUserLeave(ctx context.Context, groupID string, userID string, username string) error {
	panic("unimplemented")
}

// GroupUsersAdd implements runtime.NakamaModule.
func (m *MockModule) GroupUsersAdd(ctx context.Context, callerID string, groupID string, userIDs []string) error {
	panic("unimplemented")
}

// GroupUsersBan implements runtime.NakamaModule.
func (m *MockModule) GroupUsersBan(ctx context.Context, callerID string, groupID string, userIDs []string) error {
	panic("unimplemented")
}

// GroupUsersDemote implements runtime.NakamaModule.
func (m *MockModule) GroupUsersDemote(ctx context.Context, callerID string, groupID string, userIDs []string) error {
	panic("unimplemented")
}

// GroupUsersKick implements runtime.NakamaModule.
func (m *MockModule) GroupUsersKick(ctx context.Context, callerID string, groupID string, userIDs []string) error {
	panic("unimplemented")
}

// GroupUsersList implements runtime.NakamaModule.
func (m *MockModule) GroupUsersList(ctx context.Context, id string, limit int, state *int, cursor string) ([]*api.GroupUserList_GroupUser, string, error) {
	panic("unimplemented")
}

// GroupUsersPromote implements runtime.NakamaModule.
func (m *MockModule) GroupUsersPromote(ctx context.Context, callerID string, groupID string, userIDs []string) error {
	panic("unimplemented")
}

// GroupsGetId implements runtime.NakamaModule.
func (m *MockModule) GroupsGetId(ctx context.Context, groupIDs []string) ([]*api.Group, error) {
	panic("unimplemented")
}

// GroupsGetRandom implements runtime.NakamaModule.
func (m *MockModule) GroupsGetRandom(ctx context.Context, count int) ([]*api.Group, error) {
	panic("unimplemented")
}

// GroupsList implements runtime.NakamaModule.
func (m *MockModule) GroupsList(ctx context.Context, name string, langTag string, members *int, open *bool, limit int, cursor string) ([]*api.Group, string, error) {
	panic("unimplemented")
}

// LeaderboardCreate implements runtime.NakamaModule.
func (m *MockModule) LeaderboardCreate(ctx context.Context, id string, authoritative bool, sortOrder string, operator string, resetSchedule string, metadata map[string]interface{}) error {
	panic("unimplemented")
}

// LeaderboardDelete implements runtime.NakamaModule.
func (m *MockModule) LeaderboardDelete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// LeaderboardList implements runtime.NakamaModule.
func (m *MockModule) LeaderboardList(limit int, cursor string) (*api.LeaderboardList, error) {
	panic("unimplemented")
}

// LeaderboardRecordDelete implements runtime.NakamaModule.
func (m *MockModule) LeaderboardRecordDelete(ctx context.Context, id string, ownerID string) error {
	panic("unimplemented")
}

// LeaderboardRecordWrite implements runtime.NakamaModule.
func (m *MockModule) LeaderboardRecordWrite(ctx context.Context, id string, ownerID string, username string, score int64, subscore int64, metadata map[string]interface{}, overrideOperator *int) (*api.LeaderboardRecord, error) {
	panic("unimplemented")
}

// LeaderboardRecordsHaystack implements runtime.NakamaModule.
func (m *MockModule) LeaderboardRecordsHaystack(ctx context.Context, id string, ownerID string, limit int, cursor string, expiry int64) (*api.LeaderboardRecordList, error) {
	panic("unimplemented")
}

// LeaderboardRecordsList implements runtime.NakamaModule.
func (m *MockModule) LeaderboardRecordsList(ctx context.Context, id string, ownerIDs []string, limit int, cursor string, expiry int64) (records []*api.LeaderboardRecord, ownerRecords []*api.LeaderboardRecord, nextCursor string, prevCursor string, err error) {
	panic("unimplemented")
}

// LeaderboardRecordsListCursorFromRank implements runtime.NakamaModule.
func (m *MockModule) LeaderboardRecordsListCursorFromRank(id string, rank int64, overrideExpiry int64) (string, error) {
	panic("unimplemented")
}

// LeaderboardsGetId implements runtime.NakamaModule.
func (m *MockModule) LeaderboardsGetId(ctx context.Context, ids []string) ([]*api.Leaderboard, error) {
	panic("unimplemented")
}

// LinkApple implements runtime.NakamaModule.
func (m *MockModule) LinkApple(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// LinkCustom implements runtime.NakamaModule.
func (m *MockModule) LinkCustom(ctx context.Context, userID string, customID string) error {
	panic("unimplemented")
}

// LinkDevice implements runtime.NakamaModule.
func (m *MockModule) LinkDevice(ctx context.Context, userID string, deviceID string) error {
	panic("unimplemented")
}

// LinkEmail implements runtime.NakamaModule.
func (m *MockModule) LinkEmail(ctx context.Context, userID string, email string, password string) error {
	panic("unimplemented")
}

// LinkFacebook implements runtime.NakamaModule.
func (m *MockModule) LinkFacebook(ctx context.Context, userID string, username string, token string, importFriends bool) error {
	panic("unimplemented")
}

// LinkFacebookInstantGame implements runtime.NakamaModule.
func (m *MockModule) LinkFacebookInstantGame(ctx context.Context, userID string, signedPlayerInfo string) error {
	panic("unimplemented")
}

// LinkGameCenter implements runtime.NakamaModule.
func (m *MockModule) LinkGameCenter(ctx context.Context, userID string, playerID string, bundleID string, timestamp int64, salt string, signature string, publicKeyUrl string) error {
	panic("unimplemented")
}

// LinkGoogle implements runtime.NakamaModule.
func (m *MockModule) LinkGoogle(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// LinkSteam implements runtime.NakamaModule.
func (m *MockModule) LinkSteam(ctx context.Context, userID string, username string, token string, importFriends bool) error {
	panic("unimplemented")
}

// MatchCreate implements runtime.NakamaModule.
func (m *MockModule) MatchCreate(ctx context.Context, module string, params map[string]interface{}) (string, error) {
	panic("unimplemented")
}

// MatchGet implements runtime.NakamaModule.
func (m *MockModule) MatchGet(ctx context.Context, id string) (*api.Match, error) {
	panic("unimplemented")
}

// MatchList implements runtime.NakamaModule.
func (m *MockModule) MatchList(ctx context.Context, limit int, authoritative bool, label string, minSize *int, maxSize *int, query string) ([]*api.Match, error) {
	panic("unimplemented")
}

// MatchSignal implements runtime.NakamaModule.
func (m *MockModule) MatchSignal(ctx context.Context, id string, data string) (string, error) {
	panic("unimplemented")
}

// MetricsCounterAdd implements runtime.NakamaModule.
func (m *MockModule) MetricsCounterAdd(name string, tags map[string]string, delta int64) {
	panic("unimplemented")
}

// MetricsGaugeSet implements runtime.NakamaModule.
func (m *MockModule) MetricsGaugeSet(name string, tags map[string]string, value float64) {
	panic("unimplemented")
}

// MetricsTimerRecord implements runtime.NakamaModule.
func (m *MockModule) MetricsTimerRecord(name string, tags map[string]string, value time.Duration) {
	panic("unimplemented")
}

// MultiUpdate implements runtime.NakamaModule.
func (m *MockModule) MultiUpdate(ctx context.Context, accountUpdates []*runtime.AccountUpdate, storageWrites []*runtime.StorageWrite, walletUpdates []*runtime.WalletUpdate, updateLedger bool) ([]*api.StorageObjectAck, []*runtime.WalletUpdateResult, error) {
	panic("unimplemented")
}

// NotificationSend implements runtime.NakamaModule.
func (m *MockModule) NotificationSend(ctx context.Context, userID string, subject string, content map[string]interface{}, code int, sender string, persistent bool) error {
	panic("unimplemented")
}

// NotificationSendAll implements runtime.NakamaModule.
func (m *MockModule) NotificationSendAll(ctx context.Context, subject string, content map[string]interface{}, code int, persistent bool) error {
	panic("unimplemented")
}

// NotificationsDelete implements runtime.NakamaModule.
func (m *MockModule) NotificationsDelete(ctx context.Context, notifications []*runtime.NotificationDelete) error {
	panic("unimplemented")
}

// NotificationsSend implements runtime.NakamaModule.
func (m *MockModule) NotificationsSend(ctx context.Context, notifications []*runtime.NotificationSend) error {
	panic("unimplemented")
}

// PurchaseGetByTransactionId implements runtime.NakamaModule.
func (m *MockModule) PurchaseGetByTransactionId(ctx context.Context, transactionID string) (*api.ValidatedPurchase, error) {
	panic("unimplemented")
}

// PurchaseValidateApple implements runtime.NakamaModule.
func (m *MockModule) PurchaseValidateApple(ctx context.Context, userID string, receipt string, persist bool, passwordOverride ...string) (*api.ValidatePurchaseResponse, error) {
	panic("unimplemented")
}

// PurchaseValidateFacebookInstant implements runtime.NakamaModule.
func (m *MockModule) PurchaseValidateFacebookInstant(ctx context.Context, userID string, signedRequest string, persist bool) (*api.ValidatePurchaseResponse, error) {
	panic("unimplemented")
}

// PurchaseValidateGoogle implements runtime.NakamaModule.
func (m *MockModule) PurchaseValidateGoogle(ctx context.Context, userID string, receipt string, persist bool, overrides ...struct {
	ClientEmail string
	PrivateKey  string
}) (*api.ValidatePurchaseResponse, error) {
	panic("unimplemented")
}

// PurchaseValidateHuawei implements runtime.NakamaModule.
func (m *MockModule) PurchaseValidateHuawei(ctx context.Context, userID string, signature string, inAppPurchaseData string, persist bool) (*api.ValidatePurchaseResponse, error) {
	panic("unimplemented")
}

// PurchasesList implements runtime.NakamaModule.
func (m *MockModule) PurchasesList(ctx context.Context, userID string, limit int, cursor string) (*api.PurchaseList, error) {
	panic("unimplemented")
}

// ReadFile implements runtime.NakamaModule.
func (m *MockModule) ReadFile(path string) (*os.File, error) {
	panic("unimplemented")
}

// SessionDisconnect implements runtime.NakamaModule.
func (m *MockModule) SessionDisconnect(ctx context.Context, sessionID string, reason ...runtime.PresenceReason) error {
	panic("unimplemented")
}

// SessionLogout implements runtime.NakamaModule.
func (m *MockModule) SessionLogout(userID string, token string, refreshToken string) error {
	panic("unimplemented")
}

// StorageDelete implements runtime.NakamaModule.
func (m *MockModule) StorageDelete(ctx context.Context, deletes []*runtime.StorageDelete) error {
	panic("unimplemented")
}

// StorageIndexList implements runtime.NakamaModule.
func (m *MockModule) StorageIndexList(ctx context.Context, callerID string, indexName string, query string, limit int) (*api.StorageObjects, error) {
	panic("unimplemented")
}

// StorageList implements runtime.NakamaModule.
func (m *MockModule) StorageList(ctx context.Context, callerID string, userID string, collection string, limit int, cursor string) ([]*api.StorageObject, string, error) {
	panic("unimplemented")
}

// StorageRead implements runtime.NakamaModule.
func (m *MockModule) StorageRead(ctx context.Context, reads []*runtime.StorageRead) ([]*api.StorageObject, error) {
	panic("unimplemented")
}

// StorageWrite implements runtime.NakamaModule.
func (m *MockModule) StorageWrite(ctx context.Context, writes []*runtime.StorageWrite) ([]*api.StorageObjectAck, error) {
	panic("unimplemented")
}

// StreamClose implements runtime.NakamaModule.
func (m *MockModule) StreamClose(mode uint8, subject string, subcontext string, label string) error {
	panic("unimplemented")
}

// StreamCount implements runtime.NakamaModule.
func (m *MockModule) StreamCount(mode uint8, subject string, subcontext string, label string) (int, error) {
	panic("unimplemented")
}

// StreamSend implements runtime.NakamaModule.
func (m *MockModule) StreamSend(mode uint8, subject string, subcontext string, label string, data string, presences []runtime.Presence, reliable bool) error {
	panic("unimplemented")
}

// StreamSendRaw implements runtime.NakamaModule.
func (m *MockModule) StreamSendRaw(mode uint8, subject string, subcontext string, label string, msg *rtapi.Envelope, presences []runtime.Presence, reliable bool) error {
	panic("unimplemented")
}

// StreamUserGet implements runtime.NakamaModule.
func (m *MockModule) StreamUserGet(mode uint8, subject string, subcontext string, label string, userID string, sessionID string) (runtime.PresenceMeta, error) {
	panic("unimplemented")
}

// StreamUserJoin implements runtime.NakamaModule.
func (m *MockModule) StreamUserJoin(mode uint8, subject string, subcontext string, label string, userID string, sessionID string, hidden bool, persistence bool, status string) (bool, error) {
	panic("unimplemented")
}

// StreamUserKick implements runtime.NakamaModule.
func (m *MockModule) StreamUserKick(mode uint8, subject string, subcontext string, label string, presence runtime.Presence) error {
	panic("unimplemented")
}

// StreamUserLeave implements runtime.NakamaModule.
func (m *MockModule) StreamUserLeave(mode uint8, subject string, subcontext string, label string, userID string, sessionID string) error {
	panic("unimplemented")
}

// StreamUserList implements runtime.NakamaModule.
func (m *MockModule) StreamUserList(mode uint8, subject string, subcontext string, label string, includeHidden bool, includeNotHidden bool) ([]runtime.Presence, error) {
	panic("unimplemented")
}

// StreamUserUpdate implements runtime.NakamaModule.
func (m *MockModule) StreamUserUpdate(mode uint8, subject string, subcontext string, label string, userID string, sessionID string, hidden bool, persistence bool, status string) error {
	panic("unimplemented")
}

// SubscriptionGetByProductId implements runtime.NakamaModule.
func (m *MockModule) SubscriptionGetByProductId(ctx context.Context, userID string, productID string) (*api.ValidatedSubscription, error) {
	panic("unimplemented")
}

// SubscriptionValidateApple implements runtime.NakamaModule.
func (m *MockModule) SubscriptionValidateApple(ctx context.Context, userID string, receipt string, persist bool, passwordOverride ...string) (*api.ValidateSubscriptionResponse, error) {
	panic("unimplemented")
}

// SubscriptionValidateGoogle implements runtime.NakamaModule.
func (m *MockModule) SubscriptionValidateGoogle(ctx context.Context, userID string, receipt string, persist bool, overrides ...struct {
	ClientEmail string
	PrivateKey  string
}) (*api.ValidateSubscriptionResponse, error) {
	panic("unimplemented")
}

// SubscriptionsList implements runtime.NakamaModule.
func (m *MockModule) SubscriptionsList(ctx context.Context, userID string, limit int, cursor string) (*api.SubscriptionList, error) {
	panic("unimplemented")
}

// TournamentAddAttempt implements runtime.NakamaModule.
func (m *MockModule) TournamentAddAttempt(ctx context.Context, id string, ownerID string, count int) error {
	panic("unimplemented")
}

// TournamentCreate implements runtime.NakamaModule.
func (m *MockModule) TournamentCreate(ctx context.Context, id string, authoritative bool, sortOrder string, operator string, resetSchedule string, metadata map[string]interface{}, title string, description string, category int, startTime int, endTime int, duration int, maxSize int, maxNumScore int, joinRequired bool) error {
	panic("unimplemented")
}

// TournamentDelete implements runtime.NakamaModule.
func (m *MockModule) TournamentDelete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// TournamentJoin implements runtime.NakamaModule.
func (m *MockModule) TournamentJoin(ctx context.Context, id string, ownerID string, username string) error {
	panic("unimplemented")
}

// TournamentList implements runtime.NakamaModule.
func (m *MockModule) TournamentList(ctx context.Context, categoryStart int, categoryEnd int, startTime int, endTime int, limit int, cursor string) (*api.TournamentList, error) {
	panic("unimplemented")
}

// TournamentRecordDelete implements runtime.NakamaModule.
func (m *MockModule) TournamentRecordDelete(ctx context.Context, id string, ownerID string) error {
	panic("unimplemented")
}

// TournamentRecordWrite implements runtime.NakamaModule.
func (m *MockModule) TournamentRecordWrite(ctx context.Context, id string, ownerID string, username string, score int64, subscore int64, metadata map[string]interface{}, operatorOverride *int) (*api.LeaderboardRecord, error) {
	panic("unimplemented")
}

// TournamentRecordsHaystack implements runtime.NakamaModule.
func (m *MockModule) TournamentRecordsHaystack(ctx context.Context, id string, ownerID string, limit int, cursor string, expiry int64) (*api.TournamentRecordList, error) {
	panic("unimplemented")
}

// TournamentRecordsList implements runtime.NakamaModule.
func (m *MockModule) TournamentRecordsList(ctx context.Context, tournamentId string, ownerIDs []string, limit int, cursor string, overrideExpiry int64) (records []*api.LeaderboardRecord, ownerRecords []*api.LeaderboardRecord, prevCursor string, nextCursor string, err error) {
	panic("unimplemented")
}

// TournamentsGetId implements runtime.NakamaModule.
func (m *MockModule) TournamentsGetId(ctx context.Context, tournamentIDs []string) ([]*api.Tournament, error) {
	panic("unimplemented")
}

// UnlinkApple implements runtime.NakamaModule.
func (m *MockModule) UnlinkApple(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// UnlinkCustom implements runtime.NakamaModule.
func (m *MockModule) UnlinkCustom(ctx context.Context, userID string, customID string) error {
	panic("unimplemented")
}

// UnlinkDevice implements runtime.NakamaModule.
func (m *MockModule) UnlinkDevice(ctx context.Context, userID string, deviceID string) error {
	panic("unimplemented")
}

// UnlinkEmail implements runtime.NakamaModule.
func (m *MockModule) UnlinkEmail(ctx context.Context, userID string, email string) error {
	panic("unimplemented")
}

// UnlinkFacebook implements runtime.NakamaModule.
func (m *MockModule) UnlinkFacebook(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// UnlinkFacebookInstantGame implements runtime.NakamaModule.
func (m *MockModule) UnlinkFacebookInstantGame(ctx context.Context, userID string, signedPlayerInfo string) error {
	panic("unimplemented")
}

// UnlinkGameCenter implements runtime.NakamaModule.
func (m *MockModule) UnlinkGameCenter(ctx context.Context, userID string, playerID string, bundleID string, timestamp int64, salt string, signature string, publicKeyUrl string) error {
	panic("unimplemented")
}

// UnlinkGoogle implements runtime.NakamaModule.
func (m *MockModule) UnlinkGoogle(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// UnlinkSteam implements runtime.NakamaModule.
func (m *MockModule) UnlinkSteam(ctx context.Context, userID string, token string) error {
	panic("unimplemented")
}

// UserGroupsList implements runtime.NakamaModule.
func (m *MockModule) UserGroupsList(ctx context.Context, userID string, limit int, state *int, cursor string) ([]*api.UserGroupList_UserGroup, string, error) {
	panic("unimplemented")
}

// UsersBanId implements runtime.NakamaModule.
func (m *MockModule) UsersBanId(ctx context.Context, userIDs []string) error {
	panic("unimplemented")
}

// UsersGetId implements runtime.NakamaModule.
func (m *MockModule) UsersGetId(ctx context.Context, userIDs []string, facebookIDs []string) ([]*api.User, error) {
	panic("unimplemented")
}

// UsersGetRandom implements runtime.NakamaModule.
func (m *MockModule) UsersGetRandom(ctx context.Context, count int) ([]*api.User, error) {
	panic("unimplemented")
}

// UsersGetUsername implements runtime.NakamaModule.
func (m *MockModule) UsersGetUsername(ctx context.Context, usernames []string) ([]*api.User, error) {
	panic("unimplemented")
}

// UsersUnbanId implements runtime.NakamaModule.
func (m *MockModule) UsersUnbanId(ctx context.Context, userIDs []string) error {
	panic("unimplemented")
}

// WalletLedgerList implements runtime.NakamaModule.
func (m *MockModule) WalletLedgerList(ctx context.Context, userID string, limit int, cursor string) ([]runtime.WalletLedgerItem, string, error) {
	panic("unimplemented")
}

// WalletLedgerUpdate implements runtime.NakamaModule.
func (m *MockModule) WalletLedgerUpdate(ctx context.Context, itemID string, metadata map[string]interface{}) (runtime.WalletLedgerItem, error) {
	panic("unimplemented")
}

// WalletUpdate implements runtime.NakamaModule.
func (m *MockModule) WalletUpdate(ctx context.Context, userID string, changeset map[string]int64, metadata map[string]interface{}, updateLedger bool) (updated map[string]int64, previous map[string]int64, err error) {
	panic("unimplemented")
}

// WalletsUpdate implements runtime.NakamaModule.
func (m *MockModule) WalletsUpdate(ctx context.Context, updates []*runtime.WalletUpdate, updateLedger bool) ([]*runtime.WalletUpdateResult, error) {
	panic("unimplemented")
}
