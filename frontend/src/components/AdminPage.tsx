import React from 'react';
import { Container, Box, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Button } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

interface Contact {
  ID: number;
  Name: string;
  Email: string;
  Message: string;
  CreatedAt: string;
}

const AdminPage = () => {
  const navigate = useNavigate();
  const [contacts, setContacts] = React.useState<Contact[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/');
      return;
    }

    const fetchContacts = async () => {
      try {
        setLoading(true);
        const response = await axios.get<Contact[]>(`${import.meta.env.VITE_API_BASE_ADDRESS}/api/contacts`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        setContacts(response.data ?? []);
      } catch (error) {
        console.error('Error fetching contacts:', error as Error);
        setError((error as Error).message);
      } finally {
        setLoading(false);
      }
    };

    fetchContacts();
  }, [navigate]);

  return (
    <Container maxWidth="lg" sx={{ mt: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1">
          Admin Dashboard
        </Typography>
        <Button
          variant="outlined"
          color="secondary"
          onClick={() => {
            localStorage.removeItem('token');
            localStorage.removeItem('adminPhoneNumber');
            navigate('/');
          }}
        >
          Logout
        </Button>
      </Box>

      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '300px' }}>
          Loading...
        </Box>
      ) : error ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '300px' }}>
          {error}
        </Box>
      ) : (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Message</TableCell>
                <TableCell>Created At</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {contacts.map((contact: Contact) => (
                <TableRow key={contact.ID}>
                  <TableCell>{contact.Name}</TableCell>
                  <TableCell>{contact.Email}</TableCell>
                  <TableCell>{contact.Message}</TableCell>
                  <TableCell>{new Date(contact.CreatedAt).toLocaleString()}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}
    </Container>
  );
};

export default AdminPage;
