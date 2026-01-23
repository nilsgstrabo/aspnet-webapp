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
        private readonly ILogger _logger;
        public VideoService(IConfiguration config, ILogger<VideoService> logger)
        {
            _config=config;
            _logger=logger;
        }

        public async Task UploadVideoAsync(Stream fileStream, string fileName) {
            var file= Path.Combine(_config["VIDEO_PATH"], fileName);
            
            byte[] buffer = new byte[16384]; // Choose a suitable buffer size

            int bytesRead;
            

            using (var targetStream = new FileStream(file, new FileStreamOptions{Access=FileAccess.ReadWrite, Mode=FileMode.Create, BufferSize=1024*16}))
            {
                while ((bytesRead = await fileStream.ReadAsync(buffer)) > 0)
                {
                    _logger.LogInformation("read {bytesRead} bytes from stream", bytesRead);
                    await targetStream.WriteAsync(buffer, 0, bytesRead);
                }
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