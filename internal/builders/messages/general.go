package messages

func BuildStartMessage() string {
	return "👋 <b>Добро пожаловать в Watchlist Bot!</b>\n\n" +
		"На данный момент бот находится в разработке, но основные функции уже доступны. 🚀\n\n" +
		"Введите /help, чтобы увидеть список доступных команд."
}

func BuildHelpMessage() string {
	return "📖 <b>Список доступных команд:</b>\n\n" +
		"🎬 /start - начать пользоваться ботом\n" +
		"❓ /help - список доступных команд\n" +
		"👤 /profile - информация о вашем аккаунте\n" +
		"📚 /collections - ваши коллекции фильмов\n" +
		"⚙️ /settings - настройки\n" +
		"💬 /feedback - обратная связь\n" +
		"🚪 /logout - выйти из системы\n\n" +
		"Если у вас есть вопросы, напишите администратору. 😊"
}

func BuildFeedbackMessage() string {
	return "📝 <b>Оставьте ваш фидбек</b>\n\n" +
		"Напишите ваше мнение, предложения или идеи, и мы обязательно их учтем! 😊\n\n" +
		"Выберите категорию"
}