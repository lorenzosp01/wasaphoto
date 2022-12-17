import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

if (localStorage.getItem('token')) {
	instance.defaults.headers.common['Authorization'] = `Bearer ${localStorage.getItem('token')}`;
}

export default instance;
