import axios, { AxiosError, HttpStatusCode } from 'axios';

export interface ErrorResponse {
    error: boolean;
    statusCode: number;
}

const handleAxiosErrors = (error: unknown): ErrorResponse => {
    if (axios.isAxiosError(error)) {
        const axiosError = error as AxiosError;
        if (axiosError.response) {
            return {
                error: true,
                statusCode: 
                    axiosError.response.status > HttpStatusCode.InternalServerError
                    ? HttpStatusCode.InternalServerError
                    : axiosError.response.status
            }
        }
    }
    return { error: true, statusCode:  HttpStatusCode.InternalServerError}
}

export default handleAxiosErrors;