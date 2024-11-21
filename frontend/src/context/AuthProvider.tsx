// useAuth.ts
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { createContext, ReactNode, useCallback } from "react";
import { getToken, removeToken, storeToken } from "../utils/auth";
import axios from "axios";
import { useNavigate, useLocation } from 'react-router-dom';
import { AuthContextType, UserProfile } from "../types/auth";


export const AuthContext = createContext<AuthContextType | undefined>(undefined);

const fetchUser = async () => {
    const token = getToken();
    if (!token) return null;

    const response = await axios.get(`api/api/user`, {
        headers: {
            Authorization: `Bearer ${token}`,
        },
    });

    if (response.status !== 200) {
        throw new Error('Failed to fetch user');
    }
    const userData = response.data.user
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

export const AuthProvider = ({ children }: { children: ReactNode }) => {
    const queryClient = useQueryClient();
    const navigate = useNavigate();
    const location = useLocation();

    const { data: user, isLoading } = useQuery({
        queryKey: ['user'],
        queryFn: fetchUser,
        retry: false,
        staleTime: Infinity,
    });

    const login = useCallback((userData: UserProfile, token: string) => {
        queryClient.setQueryData(['user'], userData);
        storeToken(token);
    }, [queryClient]);

    const logout = useCallback(() => {
        removeToken();
        queryClient.setQueryData(['user'], null);
        navigate('/login');
    }, [queryClient, navigate]);

    const handleAuthCallback = useCallback(() => {
        const params = new URLSearchParams(location.search);
        const encodedData = params.get('data');
        const error = params.get('error');

        if (error) {
            console.error('Authentication error:', error);
            navigate('/login', { state: { error } });
            return;
        }

        if (encodedData) {
            try {
                const decodedData = atob(encodedData);
                const data = JSON.parse(decodedData);
                
                if (data.user && data.token) {
                    login(data.user, data.token);
                    navigate('/dashboard');
                } else {
                    throw new Error('Invalid data structure');
                }
            } catch (error) {
                console.error('Failed to process authentication data:', error);
                navigate('/login', { state: { error: 'Failed to process authentication data' } });
            }
        } else {
            navigate('/login', { state: { error: 'No authentication data received' } });
        }
    }, [location, login, navigate]);

    return (
        <AuthContext.Provider value={{ user: user ?? null, login, logout, isLoading, handleAuthCallback }}>
            {children}
        </AuthContext.Provider>
    );
};