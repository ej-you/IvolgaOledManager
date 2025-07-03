package buttons

import (
	"errors"
	"log"

	"IvolgaOledManager/internal/app/entity"
	"IvolgaOledManager/internal/pkg/errlog"
)

// screenMessage sets "message" app-status and update render.
// But if delete button is selected it deletes all messages with selected level.
func (b *Buttons) screenMessage() {
	levelMenu := b.store.Menu.GetLevel()
	selectedItem := levelMenu.Items[levelMenu.SelectedItem]
	// get selected message ID from level menu (or level if item is delete button)
	selectedMsgValue, ok := selectedItem.Value.(string)
	if !ok {
		errlog.Print(errors.New("message id is not string"))
		return
	}

	// check if delete button is selected
	if selectedItem.IsDeleteButton() {
		b.deleteWithLevel(selectedMsgValue)
		return
	}

	msg := &entity.Message{ID: selectedMsgValue}
	// get message by ID
	err := b.msgRepoDB.GetByID(msg)
	if err != nil {
		errlog.Print(err)
		return
	}

	msg.Format()
	b.store.Message.Set(msg)
	b.store.App.SetMessage()
	b.render <- struct{}{}
}

// deleteWithLevel deletes all messages with selected level and display previous menu.
func (b *Buttons) deleteWithLevel(level string) {
	err := b.msgRepoDB.DeleteAllWithLevel(level)
	if err != nil {
		errlog.Print(err)
		return
	}
	log.Printf("All messages with level %s was deleted successfully", level)
	b.screenMenuLogs()
}

// deleteMessage deletes opened message and display previous menu.
func (b *Buttons) deleteMessage() {
	msg := b.store.Message.Get()
	if !msg.IsDeleteButton() {
		return
	}
	err := b.msgRepoDB.DeleteByID(msg)
	if err != nil {
		errlog.Print(err)
		return
	}
	log.Printf("Message %q (level %s) was deleted successfully", msg.Header, msg.Level)
	b.screenMenuLevel()
}

// messageScrollUp scrolls message text up for one line.
func (b *Buttons) messageScrollUp(msg *entity.Message) {
	// scroll up message
	msg.ScrollUp()
	// update render with new message view
	b.render <- struct{}{}
}

// messageScrollDown scrolls message text down for one line.
func (b *Buttons) messageScrollDown(msg *entity.Message) {
	// scroll up message
	msg.ScrollDown()
	// update render with new message view
	b.render <- struct{}{}
}
