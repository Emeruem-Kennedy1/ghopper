import axios from "axios";
import { UserProfile } from "../types/auth";
import { getToken } from "../utils/auth";

const fetchUser = async () => {
  const token = getToken();
  if (!token) return null;

  const response = await axios.get(`api/api/user`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  if (response.status !== 200) {
    throw new Error("Failed to fetch user");
  }
  const userData = response.data.user;
  const user: UserProfile = {
    id: userData.id,
    display_name: userData.display_name,
    email: userData.email,
    uri: userData.uri,
    country: userData.country,
    image: userData.profile_image,
  };
  return user;
};

export default fetchUser;
