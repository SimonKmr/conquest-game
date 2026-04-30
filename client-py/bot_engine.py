
class BotEngine:
    def reason(self, updates, settings):
        for update in updates:
            print(update["Id"])
            print(update["X"])
            print(update["Y"])

        x = 100
        y = 200

        obj = {"Id": settings.id, "X": x, "Y": y}
        return obj