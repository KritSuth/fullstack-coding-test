import { groupByDepartment } from "../userService";
import { User } from "../types";

const mockUsers: User[] = [
  {
    id: 1,
    firstName: "Alice",
    lastName: "Smith",
    age: 30,
    gender: "female",
    hair: { color: "Black", type: "straight" },
    address: { postalCode: "10001" },
    company: { department: "Engineering" },
  },
  {
    id: 2,
    firstName: "Bob",
    lastName: "Jones",
    age: 45,
    gender: "male",
    hair: { color: "Blond", type: "curly" },
    address: { postalCode: "10002" },
    company: { department: "Engineering" },
  },
  {
    id: 3,
    firstName: "Carol",
    lastName: "White",
    age: 28,
    gender: "female",
    hair: { color: "Black", type: "wavy" },
    address: { postalCode: "20001" },
    company: { department: "Marketing" },
  },
];

describe("groupByDepartment", () => {
  const result = groupByDepartment(mockUsers);

  it("produces correct department keys", () => {
    expect(Object.keys(result).sort()).toEqual(["Engineering", "Marketing"]);
  });

  it("counts genders correctly", () => {
    expect(result["Engineering"].male).toBe(1);
    expect(result["Engineering"].female).toBe(1);
    expect(result["Marketing"].female).toBe(1);
    expect(result["Marketing"].male).toBe(0);
  });

  it("computes age range correctly", () => {
    expect(result["Engineering"].ageRange).toBe("30-45");
    expect(result["Marketing"].ageRange).toBe("28-28");
  });

  it("aggregates hair colors", () => {
    expect(result["Engineering"].hair["Black"]).toBe(1);
    expect(result["Engineering"].hair["Blond"]).toBe(1);
    expect(result["Marketing"].hair["Black"]).toBe(1);
  });

  it("maps address users", () => {
    expect(result["Engineering"].addressUser["AliceSmith"]).toBe("10001");
    expect(result["Engineering"].addressUser["BobJones"]).toBe("10002");
    expect(result["Marketing"].addressUser["CarolWhite"]).toBe("20001");
  });
});