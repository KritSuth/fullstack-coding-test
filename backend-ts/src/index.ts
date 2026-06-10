import express, { Request, Response } from "express";
import { fetchAndGroupUsers } from "./userService";
 
const app = express();
const PORT = 3000;
 
app.get("/users/grouped", async (_req: Request, res: Response) => {
  try {
    const grouped = await fetchAndGroupUsers();
    res.json(grouped);
  } catch (err) {
    res.status(500).json({ error: "Failed to fetch users" });
  }
});
 
app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});
 