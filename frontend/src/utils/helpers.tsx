export const formatUser = (responseData: any) =>
  responseData && {
    id: responseData.id,
    email: responseData.email,
    role: responseData.role,
    username: responseData.username,
  };