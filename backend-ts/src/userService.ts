import { User, GroupedByDepartment, DepartmentSummary } from "./types";
 
const API_URL = "https://dummyjson.com/users?limit=0"; // limit=0 fetches all users
 
export async function fetchAndGroupUsers(): Promise<GroupedByDepartment> {
  const response = await fetch(API_URL);
  if (!response.ok) throw new Error(`API error: ${response.status}`);
 
  const { users } = await response.json() as { users: User[] };
  return groupByDepartment(users);
}
 
export function groupByDepartment(users: User[]): GroupedByDepartment {
  const departmentMap = new Map<string, User[]>();
 
  // Group users by department in a single pass — O(n)
  for (const user of users) {
    const dept = user.company.department;
    if (!departmentMap.has(dept)) departmentMap.set(dept, []);
    departmentMap.get(dept)!.push(user);
  }
 
  const result: GroupedByDepartment = {};
 
  for (const [dept, deptUsers] of departmentMap) {
    result[dept] = summarizeDepartment(deptUsers);
  }
 
  return result;
}
 
function summarizeDepartment(users: User[]): DepartmentSummary {
  let male = 0;
  let female = 0;
  let minAge = Infinity;
  let maxAge = -Infinity;
  const hair: Record<string, number> = {};
  const addressUser: Record<string, string> = {};
 
  for (const user of users) {
    // Gender count
    if (user.gender === "male") male++;
    else female++;
 
    // Age range
    if (user.age < minAge) minAge = user.age;
    if (user.age > maxAge) maxAge = user.age;
 
    // Hair color count
    const color = user.hair.color;
    hair[color] = (hair[color] ?? 0) + 1;
 
    // Address map: firstNamelastName → postalCode
    const key = `${user.firstName}${user.lastName}`;
    addressUser[key] = user.address.postalCode;
  }
 
  return {
    male,
    female,
    ageRange: `${minAge}-${maxAge}`,
    hair,
    addressUser,
  };
}