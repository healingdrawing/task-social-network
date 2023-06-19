import axios from 'axios';
import { User } from './types';


export const signupUser = async (userData: User): Promise<User> => {
  const response = await axios.post('/api/signup', userData);
  return response.data;
};
