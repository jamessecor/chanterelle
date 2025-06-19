import React from 'react';
import { Container, Box, Typography, CssBaseline } from '@mui/material';
import ContactForm from './components/ContactForm';
import { ThemeProvider, createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ flexGrow: 1 }}>
        <Container maxWidth="lg">
          <Box sx={{ my: 4 }}>
            <Typography variant="h3" component="h1" gutterBottom align="center">
              Welcome to Our Band
            </Typography>
            <Typography variant="h6" component="h2" paragraph align="center">
              We're excited to hear from you!
            </Typography>
            <ContactForm />
          </Box>
        </Container>
      </Box>
    </ThemeProvider>
  );
}

export default App;
