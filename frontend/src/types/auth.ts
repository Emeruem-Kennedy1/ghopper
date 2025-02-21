export type UserProfile = {
    id: string;
    display_name: string;
    email: string;
    uri: string;
    country: string;
    image: string;
};

export type AuthContextType = {
    user: UserProfile | null;
    login: (userData: UserProfile, token: string) => void;
    logout: () => void;
    deleteAccount: () => Promise<void>; // Add this line
    isLoading: boolean;
    handleAuthCallback: () => void;
};