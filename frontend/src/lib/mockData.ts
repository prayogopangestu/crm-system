export interface Contact {
  id: string;
  name: string;
  email: string;
  company: string;
  role: string;
  status: 'Negosiasi' | 'Menang' | 'Prospek Awal' | 'Proposal' | 'Kalah' | 'Kualifikasi';
  lastContacted: string;
  initials: string;
  avatarUrl?: string;
}

export interface Deal {
  id: string;
  title: string;
  company: string;
  value: number;
  priority: 'High' | 'Medium' | 'Low';
  stage: 'lead' | 'contacted' | 'meeting' | 'negotiation' | 'won';
  assignee: {
    name: string;
    avatarUrl: string;
  };
}

export interface Task {
  id: string;
  title: string;
  company: string;
  time: string;
  type: 'Meeting' | 'Call' | 'Proposal' | 'Other';
  status: 'overdue' | 'today' | 'upcoming';
  completed: boolean;
}

export interface Activity {
  id: string;
  user: string;
  action: string;
  target: string;
  time: string;
  isHighlight?: boolean;
}

export interface TeamMember {
  id: string;
  name: string;
  email: string;
  role: 'Admin' | 'Staf Sales';
  status: 'Aktif' | 'Offline';
  initials: string;
}

export interface PerformanceGoal {
  month: string;
  goal: number;
  actual: number;
  status: string;
  percentage: number;
}

// Initial Mock Data
export const initialContacts: Contact[] = [
  { id: '1', name: 'Budi Wijaya', email: 'budi.w@telkomsel.co.id', company: 'PT Telkomsel', role: 'VP Sales', status: 'Negosiasi', lastContacted: 'Hari ini, 10:30', initials: 'BW' },
  { id: '2', name: 'Anita Sari', email: 'anita.sari@gojek.com', company: 'Gojek Indonesia', role: 'Procurement Manager', status: 'Menang', lastContacted: 'Kemarin', initials: 'AS' },
  { id: '3', name: 'Dian Pratama', email: 'd.pratama@bca.co.id', company: 'Bank BCA', role: 'Head of IT Infrastructure', status: 'Prospek Awal', lastContacted: '12 Okt 2023', initials: 'DP' },
  { id: '4', name: 'Hendra Kusuma', email: 'hendra@astra.co.id', company: 'PT Astra International', role: 'Supply Chain Director', status: 'Proposal', lastContacted: '10 Okt 2023', initials: 'HK', avatarUrl: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=100' },
  { id: '5', name: 'Rina Novita', email: 'rnovita@unilever.com', company: 'Unilever Indonesia', role: 'Marketing Manager', status: 'Kalah', lastContacted: '05 Okt 2023', initials: 'RN' },
  { id: '6', name: 'Fajar Ardiansyah', email: 'fajar.a@indofood.co.id', company: 'PT Indofood CBP', role: 'Operations Head', status: 'Kualifikasi', lastContacted: '02 Okt 2023', initials: 'FA' },
  { id: '7', name: 'Siti Rahma', email: 'siti.rahma@pertamina.com', company: 'Pertamina', role: 'Vendor Management', status: 'Negosiasi', lastContacted: '28 Sep 2023', initials: 'SR', avatarUrl: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=crop&q=80&w=100' },
  { id: '8', name: 'Eko Kurniawan', email: 'eko@tokopedia.com', company: 'Tokopedia', role: 'VP Engineering', status: 'Menang', lastContacted: '15 Sep 2023', initials: 'EK' }
];

export const initialDeals: Deal[] = [
  {
    id: 'd1',
    title: 'PT Maju Mundur Sentosa',
    company: 'Implementasi CRM Basic',
    value: 45000000,
    priority: 'High',
    stage: 'lead',
    assignee: {
      name: 'Andi',
      avatarUrl: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&q=80&w=100'
    }
  },
  {
    id: 'd2',
    title: 'CV Teknologi Abadi',
    company: 'Upgrade Server & Jaringan',
    value: 120000000,
    priority: 'Medium',
    stage: 'meeting',
    assignee: {
      name: 'Budi',
      avatarUrl: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&q=80&w=100'
    }
  },
  {
    id: 'd3',
    title: 'Kementerian Digital',
    company: 'Tender Pengadaan Software 2024',
    value: 850000000,
    priority: 'High',
    stage: 'negotiation',
    assignee: {
      name: 'Sarah',
      avatarUrl: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=crop&q=80&w=100'
    }
  },
  {
    id: 'd4',
    title: 'Toko Sukses Makmur',
    company: 'Langganan CRM Pro 1 Tahun',
    value: 15000000,
    priority: 'Low',
    stage: 'won',
    assignee: {
      name: 'Eko',
      avatarUrl: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=crop&q=80&w=100'
    }
  }
];

export const initialTasks: Task[] = [
  { id: 't1', title: 'Telepon Budi Wijaya', company: 'Terkait penawaran kontrak Q3', time: '14:00', type: 'Call', status: 'today', completed: false },
  { id: 't2', title: 'Follow-up Proposal', company: 'PT Maju Mundur Sentosa', time: '15:30', type: 'Proposal', status: 'today', completed: false },
  { id: 't3', title: 'Kirim Invoice Deal #1042', company: 'Implementasi IT', time: '17:00', type: 'Other', status: 'today', completed: false },
  { id: 't4', title: 'Kirim Proposal Q3', company: 'PT Garuda Abadi', time: 'Kemarin, 14:00', type: 'Proposal', status: 'overdue', completed: false },
  { id: 't5', title: 'Demo Produk CRM', company: 'Nusantara Tech', time: '10:00 - 11:30', type: 'Meeting', status: 'today', completed: true },
  { id: 't6', title: 'Makan Siang Klien', company: 'Bpk. Budi (Bank Mandiri)', time: 'Besok, 12:00', type: 'Meeting', status: 'upcoming', completed: false }
];

export const initialActivities: Activity[] = [
  { id: 'a1', user: 'Andi', action: 'menambahkan kontak baru', target: 'Budi Wijaya', time: '10 menit yang lalu' },
  { id: 'a2', user: 'Siti', action: 'memindahkan deal ke Negosiasi', target: 'CV Teknologi Abadi', time: '1 jam yang lalu' },
  { id: 'a3', user: 'Budi', action: 'menyelesaikan tugas panggilan', target: 'PT Maju Mundur Sentosa', time: '3 jam yang lalu' },
  { id: 'a4', user: 'Sistem', action: 'Deal #1042 berhasil dimenangkan', target: 'PT Telkomsel', time: 'Kemarin, 14:30', isHighlight: true },
  { id: 'a5', user: 'Sistem', action: 'membuat backup harian', target: 'Server Utama', time: 'Kemarin, 23:59' }
];

export const leaderboardData = [
  { rank: 1, name: 'Budi Santoso', role: 'Senior Sales Rep', amount: 450000000, trend: '+12%', isPositive: true, avatarUrl: 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?auto=format&fit=crop&q=80&w=100' },
  { rank: 2, name: 'Siti Rahma', role: 'Sales Executive', amount: 380000000, trend: '+5%', isPositive: true, avatarUrl: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=crop&q=80&w=100' },
  { rank: 3, name: 'Andi Wijaya', role: 'Junior Sales', amount: 210000000, trend: '-2%', isPositive: false, avatarUrl: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&q=80&w=100' }
];

export const performanceGoals: PerformanceGoal[] = [
  { month: 'Oktober', goal: 1000000000, actual: 1150000000, status: 'Tercapai (115%)', percentage: 115 },
  { month: 'September', goal: 900000000, actual: 850000000, status: 'Kurang (94%)', percentage: 94 },
  { month: 'Agustus', goal: 900000000, actual: 920000000, status: 'Tercapai (102%)', percentage: 102 }
];

export const initialTeamMembers: TeamMember[] = [
  { id: 'm1', name: 'Sarah Jenkins', email: 'sarah.j@crm-enterprise.com', role: 'Admin', status: 'Aktif', initials: 'SJ' },
  { id: 'm2', name: 'Michael Kusuma', email: 'michael.k@crm-enterprise.com', role: 'Staf Sales', status: 'Aktif', initials: 'MK' },
  { id: 'm3', name: 'Anita Larasati', email: 'anita.l@crm-enterprise.com', role: 'Staf Sales', status: 'Offline', initials: 'AL' }
];

export const pipelineStages = [
  { id: 's1', name: 'Prospek Baru', color: 'bg-primary-container' },
  { id: 's2', name: 'Kualifikasi', color: 'bg-secondary-container' },
  { id: 's3', name: 'Negosiasi', color: 'bg-tertiary-container' },
  { id: 's4', name: 'Deal Won', color: 'bg-surface-tint' }
];
