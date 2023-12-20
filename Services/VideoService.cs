using Azure.Security.KeyVault.Secrets;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Logging;
using System.Collections.Generic;
using System.Threading.Tasks;
using System.Linq;
using System.IO;
namespace aspnet_webapp.Services {

    public class VideoInfo {
        public string Id { get; set; }
        public string FileName { get; set; }
        public long SizeBytes { get; set; }
    }

    public interface IVideoService {
        public IEnumerable<VideoInfo> GetVideos();
        public Task UploadVideoAsync(Stream fileStream, string fileName);
    }

    public class VideoService : IVideoService {
        private readonly IConfiguration _config;
        public VideoService(IConfiguration config)
        {
            _config=config;
        }

        public async Task UploadVideoAsync(Stream fileStream, string fileName) {
            var file= Path.Combine(_config["VIDEO_PATH"], fileName);
            using (var targetStream = new FileStream(file, FileMode.Create))
            {
                await fileStream.CopyToAsync(targetStream);
            }
        }

        public IEnumerable<VideoInfo> GetVideos()
        {
            var videos = new List<VideoInfo>();
            var rootPath=_config["VIDEO_PATH"];
            if(rootPath?.Length==0) {
                return videos;
            }
           
            var q=new Stack<string>(new string[]{rootPath});
            
            
            while (q.Count>0)
            {
                var dir=q.Pop();
                System.IO.Directory.GetDirectories(dir).ToList().ForEach(d=>q.Push(d));
                System.IO.Directory.GetFiles(dir).Select(f=>new System.IO.FileInfo(f)).ToList().ForEach(f=>videos.Add(new VideoInfo{Id=f.FullName.Replace('/','-'),FileName=f.FullName, SizeBytes=f.Length}));
            }
            return videos;

        }
    }

}