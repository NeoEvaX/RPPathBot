package bot

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken string
	GuildID  string
	s        *discordgo.Session
)

var (
	defaultMemberPermissions int64 = discordgo.PermissionAdministrator
	commands                       = []*discordgo.ApplicationCommand{
		{
			Name:                     "setupgame2",
			Description:              "Setup a new game",
			DefaultMemberPermissions: &defaultMemberPermissions,
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, m *discordgo.InteractionCreate){
		"setupgame2": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			slog.Info("Creating Category: setup")
			newCategory, err := s.GuildChannelCreate(i.GuildID, "setup", discordgo.ChannelTypeGuildCategory)
			slog.Info("Created Category: setup")

			if err == nil {
				slog.Info("Creating Channel: channel")
				channelData := discordgo.GuildChannelCreateData{Name: "channel", Type: discordgo.ChannelTypeGuildText, ParentID: newCategory.ID}
				s.GuildChannelCreateComplex(i.GuildID, channelData)
				slog.Info("Created Channel: channel")

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Game created!",
					},
				})
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error creating game!",
					},
				})
			}
		},
	}
)

func Run() {
	// create a session
	s, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		slog.Error("Invalid bot parameters", slog.Any("err", err.Error()))
		return
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		slog.Info("Logged in as:", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// open session
	err = s.Open()
	if err != nil {
		slog.Error("Cannot open session", slog.Any("err", err.Error()))
		return
	}

	slog.Info("Adding commands...", slog.Int("amount:", len(commands)))
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			slog.Error("Cannot create '%v' command: %v", v.Name, err)
			return
		}
		registeredCommands[i] = cmd
		slog.Info("Created command: ", slog.Any("command", v.Name))
	}

	defer s.Close() // close session, after function termination

	// keep bot running until there is NO os interruption (ctrl + C)
	slog.Info("Bot running....")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	slog.Info("Press Ctrl + C to exit")
	<-stop

	slog.Info("Removing commands...")
	// // We need to fetch the commands, since deleting requires the command ID.
	// // We are doing this from the returned commands on line 375, because using
	// // this will delete all the commands, which might not be desirable, so we
	// // are deleting only the commands that we added.
	//registeredCommands, err = s.ApplicationCommands(s.State.User.ID, "")
	//if err != nil {
	//	log.Fatalf("Could not fetch registered commands: %v", err)
	//}

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
		if err != nil {
			slog.Error("Cannot delete '%v' command: %v", v.Name, err)
		} else {
			slog.Info("Deleted command", slog.String("name", v.Name))
		}
	}

	slog.Info("Bot stopped")
}

func newMessage(s *discordgo.Session, message *discordgo.MessageCreate) {
	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == s.State.User.ID {
		return
	}
}
