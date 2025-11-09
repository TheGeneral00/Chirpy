package helpers

import()

type TestUser struct {
	Email string
	Password string
}

func LoadTestUsers() []TestUser{
	return  []TestUser{
		// --- Valid, typical ---
        {"alice@example.com", "password123"},
        {"bob.smith@example.co.uk", "hunter2"},
        {"charlie99@test.io", "LetMeIn!"},
        {"dora_2025@gmail.com", "SafePass987"},
        {"eric@example.com", "p@55w0rd"},

        // --- Edge passwords ---
        {"tiny@example.com", "a"},                // too short
        {"longboi@example.com", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"}, // too long
        {"specialchars@example.com", "!@#$%^&*()_+=-[]{}"},
        {"unicode@example.com", "пароль秘密كلمةالسر"},
        {"spacey@example.com", "   padded   "},

        // --- Odd emails ---
        {"weird+label@example.com", "okpassword"},
        {"upperCASE@Example.COM", "mixedCase1"},
        {"emoji😊@example.com", "emojiPass"},
        {"noatsymbol.com", "bademail1"},
        {"@@doubleats@@", "bademail2"},

        // --- Potential injections / attacks ---
        {"sqlinjection@example.com' OR '1'='1", "sqltest"},
        {"xss<script>@example.com", "<script>alert('xss')</script>"},
        {"newline@ex\nample.com", "newlinepass"},
        {"tab@exa\tmple.com", "tabpass"},
        {"quote\"@example.com", "quote'pass"},

        // --- Whitespace / invisible ---
        {" leading@example.com", "spaceStart"},
        {"trailing@example.com ", "spaceEnd"},
        {"invisible\u200b@example.com", "zeroWidth"},
	}
}
