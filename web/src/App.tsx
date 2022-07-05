import { useState } from 'react'
import { AppBar, Toolbar, Typography, Button, Container, Box, CircularProgress, Divider, Dialog, DialogTitle, DialogContent } from '@mui/material'
import ky from 'ky';

function App() {
  const [imgSrc, setImgSrc] = useState("");
  const [loading, setLoading] = useState(false);
  const [logs, setLogs] = useState("");
  const [showLogs, setShowLogs] = useState(false);
  async function fetchImage() {
    setLoading(true);
    const blob = await ky(`${import.meta.env.VITE_BASE_API}imageAPI`, { cache: 'no-store', timeout: 60000 }).blob()
    setLoading(false);
    setImgSrc(URL.createObjectURL(blob))
  }
  async function setRate(rate: number) {
    await ky(`${import.meta.env.VITE_BASE_API}chance`, { cache: 'no-store', method: 'post', json: { chance: rate, timeout: 60000 } })
  }
  async function setPattern(pattern: string) {
    await ky(`${import.meta.env.VITE_BASE_API}pattern`, { cache: 'no-store', method: 'post', json: { pattern: pattern, timeout: 60000 } })
  }
  async function fetchLogs() {
    setLoading(true)
    setShowLogs(true);
    const logsResp = await ky(`${import.meta.env.VITE_BASE_API}log`, { cache: 'no-store', timeout: 60000 }).json();
    setLogs((logsResp as any).contents);
    setLoading(false);
  }
  return (
    <>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Software Patterns - AWS Interview
          </Typography>
          <Button color="inherit" onClick={() => fetchLogs()}>View Logs</Button>
        </Toolbar>
      </AppBar>
      <Container sx={{ mt: 2 }}>
        <Typography variant="h6">Set Failure Rate</Typography>
        <Box>
          <Button onClick={() => setRate(0.00000001)}>0%</Button>
          <Button onClick={() => setRate(25)}>25%</Button>
          <Button onClick={() => setRate(50)}>50%</Button>
          <Button onClick={() => setRate(75)}>75%</Button>
          <Button onClick={() => setRate(100)}>100%</Button>
        </Box>
        <Divider sx={{ my: 2 }} />
        <Typography variant="h6">Set Pattern</Typography>
        <Box>
          <Button onClick={() => setPattern("NO_PATTERN")}>No Pattern</Button>
          <Button onClick={() => setPattern("BACKOFF")}>Backoff</Button>
          <Button onClick={() => setPattern("EXPONENTIAL_BACKOFF")}>Exponential Backoff</Button>
          <Button onClick={() => setPattern("EXP_WITH_JITTER")}>Jitter + Exp Backoff</Button>
          <Button onClick={() => setPattern("CIRCUIT_BREAKER")}>Circuit Breaker</Button>
        </Box>
        <Divider sx={{ my: 2 }} />
        <Box>
          <Button onClick={fetchImage}>Fetch New Image</Button>
        </Box>
        {loading ?
          <CircularProgress size="10rem" /> :
          <Box id="image" sx={{ width: 400, height: 400, backgroundSize: "contain", backgroundImage: `url(${imgSrc})`, backgroundRepeat: "no-repeat", backgroundPosition: "center" }} />}
      </Container>
      <Dialog fullWidth open={Boolean(showLogs)} onClose={() => setShowLogs(false)}>
        <DialogTitle align="center">logs</DialogTitle>
        <DialogContent sx={{ whiteSpace: 'pre-line' }}>
          {loading ? <CircularProgress /> : <>{logs}</>}
        </DialogContent>
      </Dialog>
    </>
  )
}

export default App
